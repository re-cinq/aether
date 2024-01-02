package gcp

import (
	"context"
	"fmt"
	"strconv"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	monitoringpb "cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
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
	client    *monitoring.QueryClient
	projectID string
	cache     *gcpCache
}

type options func(*GCP)

// New returns a new instance of the GCP provider as well as a function to
// cleanup connections once done
func New(account *config.Account, cache *gcpCache, opts ...options) (g *GCP, teardown func(), err error) {
	// set any defaults here
	g = &GCP{
		projectID: account.Project,
		cache:     cache,
	}

	var clientOptions []option.ClientOption

	if account.Credentials.IsPresent() {
		credentialFile := account.Credentials.FilePaths[0]
		clientOptions = append(clientOptions, option.WithCredentialsFile(credentialFile))
	}

	// overwrite any options
	for _, opt := range opts {
		opt(g)
	}

	// This allows overwriting the default client
	// google by default trys to authenticate when initilizing
	// a client therefore if we put this before running the options
	// it would try authenticate against google regardless of overwriting the
	// client
	if g.client == nil {
		c, err := monitoring.NewQueryClient(context.TODO(), clientOptions...)
		if err != nil {
			return nil, func() {}, err
		}
		g.client = c
	}

	// teardown is used to close relevant connections
	// and cleanup
	teardown = func() {
		g.client.Close()
	}

	return g, teardown, nil
}

func (g *GCP) GetMetricsForInstances(
	ctx context.Context,
	window string,
) ([]v1.Instance, error) {
	var services []v1.Instance

	metrics, err := g.instanceMetrics(
		ctx, fmt.Sprintf(CPUQuery, g.projectID, window, window),
	)

	if err != nil {
		return services, err
	}

	for _, m := range metrics {
		metric := *m

		// Get the zone
		if zone, ok := metric.Labels().Get("zone"); ok {
			// Get the region
			if region, ok := metric.Labels().Get("region"); ok {
				// Get the instance name
				if instanceName, ok := metric.Labels().Get("name"); ok {
					if instanceID, ok := metric.Labels().Get("id"); ok {
						// Load the cacge
						cachedInstance := g.cache.Get(zone, gceService, instanceName)

						if cachedInstance != nil {
							instance := v1.NewInstance(instanceID, gcpProvider).SetService(gceService)

							if machineType, ok := metric.Labels().Get("machine_type"); ok {
								instance.SetKind(machineType)
							}

							instance.SetRegion(region)
							instance.SetZone(zone)
							instance.Metrics().Upsert(&metric)

							services = append(services, *instance)
						}
					}
				}
			}
		}
	}

	return services, nil
}

// GetCPUUtilization returns the utilization for instances and is a wrapper
func (g *GCP) GetCPUForInstances(
	ctx context.Context,
	window string,
) ([]*v1.Metric, error) {
	return g.instanceMetrics(
		ctx, fmt.Sprintf(CPUQuery, g.projectID, window, window),
	)
}

// instanceMetrics runs a query on googe cloud monitoring using MQL
// and responds with a list of metrics
func (g *GCP) instanceMetrics(
	ctx context.Context,
	query string,
) ([]*v1.Metric, error) {
	var metrics []*v1.Metric
	it := g.client.QueryTimeSeries(ctx, &monitoringpb.QueryTimeSeriesRequest{
		Name:  fmt.Sprintf("projects/%s", g.projectID),
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

		instanceID := resp.GetLabelValues()[0].GetStringValue()
		instanceName := resp.GetLabelValues()[1].GetStringValue()
		region := resp.GetLabelValues()[2].GetStringValue()
		zone := resp.GetLabelValues()[3].GetStringValue()

		var instanceType string
		if len(resp.GetLabelValues()) > 4 {
			instanceType = resp.GetLabelValues()[4].GetStringValue()
		}

		var totalCores string
		if len(resp.GetLabelValues()) > 5 {
			totalCores = resp.GetLabelValues()[5].GetStringValue()
		}

		m := v1.NewMetric("cpu")
		m.SetResourceUnit(v1.Core)
		m.SetType(v1.CPU).SetUsagePercentage(resp.GetPointData()[0].GetValues()[0].GetDoubleValue() * 100)

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
