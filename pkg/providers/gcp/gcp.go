package gcp

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	monitoringpb "cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	cache "github.com/patrickmn/go-cache"
	"github.com/re-cinq/cloud-carbon/pkg/config"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"k8s.io/klog/v2"
)

var (
	/*
	* An MQL query that will return data from Google Cloud with the
	* - Instance Name
	* - Region
	* - Zone
	* - Machine Type
	* - Reserved CPUs
	* - Utilization
	 */
	CPUQuery = `
  fetch gce_instance
  | { metric 'compute.googleapis.com/instance/cpu/utilization'
    ; metric 'compute.googleapis.com/instance/cpu/reserved_cores' }
  | outer_join 0
	| filter project_id = '%s' 
  | group_by [		
    resource.instance_id,
  	metric.instance_name,
		metadata.system.region,
		resource.zone,
		metadata.system.machine_type,
    reserved_cores: format(t_1.value.reserved_cores, '%%f')
  ], [max(t_0.value.utilization)]
  | window %s
  | within %s
	`
)

// GCP is the structure used as the provider for Google Cloud Platform
type GCP struct {
	// GCP Clients
	monitoring *monitoring.QueryClient
	instances  *compute.InstancesClient

	// Caching mechanism
	cache *cache.Cache
}

// The GCP Resource representation
type resource struct {
	// The GCP resource id
	id string

	// The region where the resource is located
	region string

	// The service the resource belongs to
	service string

	// For example spot, reserved
	lifecycle string

	// Amount of cores
	coreCount int

	// The instance kind for example n2-standard-8
	kind string

	// The name
	name string

	// When was the last time it was updated
	lastUpdated time.Time
}

type options func(*GCP)

// New returns a new instance of the GCP provider as well as a function to
// cleanup connections once done
func New(
	ctx context.Context,
	account *config.Account,
	opts ...options,
) (g *GCP, teardown func(), err error) {
	// set any defaults here
	g = &GCP{
		// TODO do we want to expire cache?
		cache: cache.New(3600*time.Minute, 3600*time.Minute),
	}

	var clientOptions []option.ClientOption

	if account.Credentials.IsPresent() {
		credentialFile := account.Credentials.FilePaths[0]
		clientOptions = append(
			clientOptions,
			option.WithCredentialsFile(credentialFile),
		)
	}

	// overwrite any options
	for _, opt := range opts {
		opt(g)
	}

	// This allows overwriting the default monitoring client
	// google by default trys to authenticate when initilizing
	// a client therefore if we put this before running the options
	// it would try authenticate against google regardless of overwriting the
	// client
	if g.monitoring == nil {
		c, err := monitoring.NewQueryClient(ctx, clientOptions...)
		if err != nil {
			return nil, func() {}, err
		}
		g.monitoring = c
	}

	// This allows overwriting the default instances client
	if g.instances == nil {
		c, err := compute.NewInstancesRESTClient(ctx, clientOptions...)
		if err != nil {
			return nil, func() {}, err
		}
		g.instances = c
	}

	// teardown is used to close relevant connections
	// and cleanup
	teardown = func() {
		g.monitoring.Close()
		g.instances.Close()
	}

	return g, teardown, nil
}

// GetMetricsForInstances retrieves all the metrics for a given instance
func (g *GCP) GetMetricsForInstances(
	ctx context.Context,
	project, window string,
) ([]v1.Instance, error) {
	var instances []v1.Instance

	metrics, err := g.instanceMetrics(
		// TODO these parameters can be cleaned up
		ctx, project, fmt.Sprintf(CPUQuery, project, window, window),
	)

	if err != nil {
		return instances, err
	}

	// TODO there seems to be duplicated logic here
	// Why not create instance whuile collecting metric instead of handeling
	// it in two steps
	for _, m := range metrics {
		metric := *m

		zone, ok := metric.Labels().Get("zone")
		if !ok {
			continue
		}

		region, ok := metric.Labels().Get("region")
		if !ok {
			continue
		}

		instanceName, ok := metric.Labels().Get("name")
		if !ok {
			continue
		}

		instanceID, ok := metric.Labels().Get("id")
		if !ok {
			continue
		}

		machineType, ok := metric.Labels().Get("machine_type")
		if !ok {
			continue
		}

		// Load the cache
		// TODO make this more explicit, im not sure why this
		// is needed as we dont use the cache anywhere
		cachedInstance, ok := g.cache.Get(cacheKey(zone, service, instanceName))
		if cachedInstance == nil && !ok {
			continue
		}

		instance := v1.NewInstance(instanceID, provider).SetService(service)

		instance.SetKind(machineType)
		instance.SetRegion(region)
		instance.SetZone(zone)
		instance.Metrics().Upsert(&metric)

		instances = append(instances, *instance)
	}

	return instances, nil
}

// instanceMetrics runs a query on googe cloud monitoring using MQL
// and responds with a list of metrics
func (g *GCP) instanceMetrics(
	ctx context.Context,
	project, query string,
) ([]*v1.Metric, error) {
	var metrics []*v1.Metric

	it := g.monitoring.QueryTimeSeries(ctx, &monitoringpb.QueryTimeSeriesRequest{
		Name:  fmt.Sprintf("projects/%s", project),
		Query: query,
	})

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		// This is dependant on the MQL query
		// label ordering
		instanceID := resp.GetLabelValues()[0].GetStringValue()
		instanceName := resp.GetLabelValues()[1].GetStringValue()
		region := resp.GetLabelValues()[2].GetStringValue()
		zone := resp.GetLabelValues()[3].GetStringValue()
		instanceType := resp.GetLabelValues()[4].GetStringValue()
		totalCores := resp.GetLabelValues()[5].GetStringValue()

		m := v1.NewMetric("cpu")
		m.SetResourceUnit(v1.Core)
		m.SetType(v1.CPU).SetUsagePercentage(
			// translate fraction to a percentage
			resp.GetPointData()[0].GetValues()[0].GetDoubleValue() * 100,
		)

		f, err := strconv.ParseFloat(totalCores, 64)
		// TODO: we should not fail here but collect errors
		if err != nil {
			klog.Errorf("failed to parse GCP metric %s", err)
			continue
		}

		m.SetUnitAmount(f)
		m.SetLabels(v1.Labels{
			"id":           instanceID,
			"name":         instanceName,
			"region":       region,
			"zone":         zone,
			"machine_type": instanceType,
		})
		metrics = append(metrics, m)
	}
	return metrics, nil
}

// Refresh fetches all the Instances
// for a project and stores metadata in order to help with
// metric collections
func (g *GCP) Refresh(ctx context.Context, project string) {
	iter := g.instances.AggregatedList(
		ctx,
		&computepb.AggregatedListInstancesRequest{
			Project: project,
		},
	)

	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			klog.Errorf("error while processes GCE instances %s", err)
			return
		}

		for _, instance := range resp.Value.Instances {
			zone := getValueFromURL(instance.GetZone())
			instanceID := strconv.FormatUint(instance.GetId(), 10)
			name := instance.GetName()

			if zone == "" {
				continue
			}

			if instance.GetStatus() == "TERMINATED" {
				// delete the entry from the cache
				g.cache.Delete(cacheKey(zone, service, name))
				continue
			}

			if instance.GetStatus() == "RUNNING" {
				// TODO potentially we do not need a custom resource to cache, maybe
				// just cache the instance
				g.cache.Set(cacheKey(zone, service, name), resource{
					region:      zone,
					service:     service,
					id:          instanceID,
					kind:        getValueFromURL(instance.GetMachineType()),
					lifecycle:   instance.GetScheduling().GetProvisioningModel(),
					name:        name,
					coreCount:   0,
					lastUpdated: time.Now().UTC(),
				}, cache.DefaultExpiration)
			}
		}
	}
}

func cacheKey(z, s, n string) string {
	return fmt.Sprintf("%s-%s-%s", z, s, n)
}

// getValueFromURL returns the last element in the url Path
// example:
// input: https://www.googleapis.com/.../machineTypes/e2-micro
// output: e2-micro
func getValueFromURL(u string) string {
	parsed, err := url.Parse(u)
	if err != nil {
		klog.Errorf("failed to parse value from %s %s", u, err)
		return ""
	}

	return path.Base(parsed.Path)
}
