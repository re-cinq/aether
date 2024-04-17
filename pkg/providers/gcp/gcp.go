package gcp

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"github.com/re-cinq/aether/pkg/config"
	"github.com/re-cinq/aether/pkg/log"
	"github.com/re-cinq/aether/pkg/providers/util"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Client is the structure used as the provider for Google Cloud Platform
type Client struct {
	// GCP Clients
	monitoring *monitoring.QueryClient
	compute    *compute.InstancesClient

	// hashmap of instances
	instancesMap map[string]*v1.Instance
}

type options func(*Client)

// New returns a new instance of the GCP provider as well as a function to
// cleanup connections once done
func New(
	ctx context.Context,
	account *config.Account,
	opts ...options,
) (c *Client, teardown func(), err error) {
	c = &Client{
		// initilize instance lookup table
		instancesMap: make(map[string]*v1.Instance),
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
		opt(c)
	}

	// This allows overwriting the default monitoring client
	// google by default trys to authenticate when initilizing
	// a client therefore if we put this before running the options
	// it would try authenticate against google regardless of overwriting the
	// client
	if c.monitoring == nil {
		mc, err := monitoring.NewQueryClient(ctx, clientOptions...)
		if err != nil {
			return nil, func() {}, err
		}
		c.monitoring = mc
	}

	// This allows overwriting the default instances client
	if c.compute == nil {
		ic, err := compute.NewInstancesRESTClient(ctx, clientOptions...)
		if err != nil {
			return nil, func() {}, err
		}
		c.compute = ic
	}

	// teardown is used to close relevant connections
	// and cleanup
	teardown = func() {
		c.monitoring.Close()
		c.compute.Close()
	}

	return c, teardown, nil
}

// GetMetricsForInstances retrieves all the metrics for a given instance
// And updates the cached instance with the metrics
func (c *Client) GetMetricsForInstances(
	ctx context.Context,
	project, window string,
) error {
	// TODO these parameters can be cleaned up
	err := c.cpuMetrics(ctx, project, fmt.Sprintf(CPUQuery, project, window, window))
	if err != nil {
		return err
	}

	err = c.memoryMetrics(ctx, project, fmt.Sprintf(MEMQuery, project, window, window))
	if err != nil {
		return err
	}

	return nil
}

// Refresh fetches all the Instances
// for a project and stores metadata in order to help with
// metric collections
func (c *Client) Refresh(ctx context.Context, project string) {
	logger := log.FromContext(ctx)

	iter := c.compute.AggregatedList(
		ctx,
		&computepb.AggregatedListInstancesRequest{
			Project: project,
		},
	)

	// instances is a slice of all valid instances
	// that will be stored as a value in the cache
	// with the key `gcp-valid-instances`
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Error("failed processesing GCE instance", "error", err)
			return
		}

		for _, instance := range resp.Value.Instances {
			instanceID := strconv.FormatUint(instance.GetId(), 10)
			name := instance.GetName()

			zone, err := getValueFromURL(instance.GetZone())
			if err != nil || zone == "" {
				logger.Error("failed to get zone", "error", err, "instanceID", instanceID)
				continue
			}

			region, err := getRegionFromZone(zone)
			if err != nil {
				logger.Error("error getting region", "error", err, "instanceID", instanceID)
				continue
			}

			kind, err := getValueFromURL(instance.GetMachineType())
			if err != nil {
				logger.Error("failed to get instance type from url")
			}

			mapInstance := &v1.Instance{
				ID:       instanceID,
				Provider: provider,
				Name:     name,
				Region:   region,
				Zone:     zone,
				Service:  service,
				Kind:     kind,
				Labels: v1.Labels{
					"Lifecycle": instance.GetScheduling().GetProvisioningModel(),
					"ID":        instanceID,
				},
			}

			if instance.GetStatus() == "TERMINATED" {
				mapInstance.Status = v1.InstanceTerminated
			}

			if instance.GetStatus() == "RUNNING" {
				mapInstance.Status = v1.InstanceRunning
			}

			// Add running instances to the cache
			key := util.Key(zone, service, name)
			c.instancesMap[key] = mapInstance
		}
	}
}

// getValueFromURL returns the last element in the url Path
// example:
// input: https://www.googleapis.com/.../machineTypes/e2-micro
// output: e2-micro
func getValueFromURL(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	return path.Base(parsed.Path), nil
}

// getRegionFromZone removes the suffix from the GCE zone
// to get the region
// input: europe-west1-a
// output: europe-west1
func getRegionFromZone(z string) (string, error) {
	x := strings.LastIndex(z, "-")
	if x == -1 {
		return "", fmt.Errorf("error: cannot get region from zone")
	}

	return z[:x], nil
}
