package gcp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"path"
	"strconv"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	cache "github.com/patrickmn/go-cache"
	"github.com/re-cinq/cloud-carbon/pkg/config"
	"github.com/re-cinq/cloud-carbon/pkg/providers/util"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// GCP is the structure used as the provider for Google Cloud Platform
type GCP struct {
	// GCP Clients
	monitoring *monitoring.QueryClient
	instances  *compute.InstancesClient

	// Caching mechanism
	cache *cache.Cache
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

	cpumetrics, err := g.instanceCPUMetrics(
		// TODO these parameters can be cleaned up
		ctx, project, fmt.Sprintf(CPUQuery, project, window, window),
	)
	if err != nil {
		return instances, err
	}

	memmetrics, err := g.instanceMemoryMetrics(
		ctx, project, fmt.Sprintf(MEMQuery, project, window, window),
	)
	if err != nil {
		return instances, err
	}

	// we use a lookup to add different metrics to the same instance
	lookup := make(map[string]*v1.Instance)

	// TODO there seems to be duplicated logic here
	// Why not create instance whuile collecting metric instead of handeling
	// it in two steps
	for _, m := range append(cpumetrics, memmetrics...) {
		metric := *m

		meta, err := getMetadata(&metric)
		if err != nil {
			continue
		}

		// Load the cache
		// TODO make this more explicit, im not sure why this
		// is needed as we dont use the cache anywhere
		// I think this is removed, and is used when the metric data doesn't have
		// the complete resource/instance information.
		cachedInstance, ok := g.cache.Get(util.CacheKey(meta.zone, service, meta.name))
		if cachedInstance == nil && !ok {
			continue
		}

		i, ok := lookup[meta.id]
		if !ok {
			i = v1.NewInstance(meta.id, provider).SetService(service)
		}

		i.SetKind(meta.machineType)
		i.SetRegion(meta.region)
		i.SetZone(meta.zone)
		i.Metrics().Upsert(&metric)

		lookup[meta.id] = i
	}

	// create list of instances
	// TODO: this seems repetitive
	for _, v := range lookup {
		instances = append(instances, *v)
	}

	return instances, nil
}

type metadata struct {
	zone, region, name, id, machineType string
}

func getMetadata(m *v1.Metric) (*metadata, error) {
	zone, ok := m.Labels().Get("zone")
	if !ok {
		return &metadata{}, errors.New("zone not found")
	}

	region, ok := m.Labels().Get("region")
	if !ok {
		return &metadata{}, errors.New("region not found")
	}

	name, ok := m.Labels().Get("name")
	if !ok {
		return &metadata{}, errors.New("instance name not found")
	}

	id, ok := m.Labels().Get("id")
	if !ok {
		return &metadata{}, errors.New("instance id not found")
	}

	machineType, ok := m.Labels().Get("machine_type")
	if !ok {
		return &metadata{}, errors.New("machine type not found")
	}
	return &metadata{
		zone:        zone,
		region:      region,
		name:        name,
		id:          id,
		machineType: machineType,
	}, nil
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
			slog.Error("failed processesing GCE instance", "error", err)
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
				g.cache.Delete(util.CacheKey(zone, service, name))
				continue
			}

			if instance.GetStatus() == "RUNNING" {
				// TODO potentially we do not need a custom resource to cache, maybe
				// just cache the instance (Kind fields are different in resource and
				// instance) - will need to consolidate
				g.cache.Set(util.CacheKey(zone, service, name), v1.Resource{
					ID:          instanceID,
					Name:        name,
					Region:      zone, // TODO: Why is region set to zone
					Service:     service,
					Kind:        getValueFromURL(instance.GetMachineType()),
					Lifecycle:   instance.GetScheduling().GetProvisioningModel(),
					VCPUCount:   0,
					LastUpdated: time.Now().UTC(),
				}, cache.DefaultExpiration)
			}
		}
	}
}

// getValueFromURL returns the last element in the url Path
// example:
// input: https://www.googleapis.com/.../machineTypes/e2-micro
// output: e2-micro
func getValueFromURL(u string) string {
	parsed, err := url.Parse(u)
	if err != nil {
		slog.Error("failed to parse value from", "input", u, "error", err)
		return ""
	}

	return path.Base(parsed.Path)
}
