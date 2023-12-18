package gcp

import (
	"context"
	"net/url"
	"path"
	"strconv"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/re-cinq/cloud-carbon/pkg/config"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"k8s.io/klog/v2"
)

type gceClient struct {
	cache  *gcpCache
	client *compute.InstancesClient
}

func newGCECLient(account config.Account) *gceClient {
	ctx := context.Background()

	var clientOptions []option.ClientOption

	if account.Credentials.IsPresent() {
		credentialFile := account.Credentials.FilePaths[0]
		clientOptions = append(clientOptions, option.WithCredentialsFile(credentialFile))
	}

	client, err := compute.NewInstancesRESTClient(ctx, clientOptions...)
	if err != nil {
		klog.Errorf("failed to create GCE client %s", err)
		return nil
	}

	return &gceClient{
		client: client,
		cache:  newGCPCache(),
	}
}

func (e *gceClient) Close() {
	e.client.Close()
}

// refresh stores all the instances for a specific project
func (e *gceClient) Refresh(project string) {
	ctx := context.Background()
	req := &computepb.AggregatedListInstancesRequest{
		Project: project,
		// See https://pkg.go.dev/cloud.google.com/go/compute/apiv1/computepb#AggregatedListInstancesRequest.
	}
	it := e.client.AggregatedList(ctx, req)
	for {
		resp, err := it.Next()
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
				e.cache.Delete(zone, gceService, name)
			} else if instance.GetStatus() == "RUNNING" {
				// If it's running add it
				e.cache.Add(newGCPResource(
					zone,
					gceService,
					instanceID,
					getValueFromURL(instance.GetMachineType()),
					instance.GetScheduling().GetProvisioningModel(),
					name,
					0, // GCE does not give us the amount of CPUs the instance has
				))
			}
		}
	}
}

// string example: https://www.googleapis.com/compute/v1/projects/cloud-carbon-project/zones/europe-north1-a/machineTypes/e2-micro
func getValueFromURL(gceURL string) string {
	parsed, err := url.Parse(gceURL)
	if err != nil {
		klog.Errorf("failed to parse value from %s %s", gceURL, err)
		return ""
	}

	return path.Base(parsed.Path)
}
