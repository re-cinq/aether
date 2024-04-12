// Contains a set of method for getting EC2 information
package amazon

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/re-cinq/aether/pkg/providers/util"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

// refresh stores all the instances for a specific region in cache
func (c *Client) Refresh(ctx context.Context, region string) error {
	// Override the region
	withRegion := func(o *ec2.Options) {
		o.Region = region
	}

	// First request
	output, err := c.ec2.DescribeInstances(ctx, buildListPaginationRequest(nil), withRegion)
	if err != nil || output == nil {
		return fmt.Errorf("failed to retrieve ec2 instances from region: %s: %s", region, err)
	}

	c.updateInstancesMap(region, output.Reservations)

	for output.NextToken != nil {
		output, err = c.ec2.DescribeInstances(ctx, buildListPaginationRequest(output.NextToken), withRegion)
		if err != nil || output == nil {
			return fmt.Errorf("failed to retrieve ec2 instances %s", err)
		}

		c.updateInstancesMap(region, output.Reservations)
	}

	return nil
}

func (c *Client) updateInstancesMap(region string, res []types.Reservation) {
	for _, r := range res {
		for index := range r.Instances {
			instance := r.Instances[index]

			id := aws.ToString(instance.InstanceId)
			key := util.Key(region, ec2Service, id)

			// Remove non-running instances
			if instance.State.Name != types.InstanceStateNameRunning {
				delete(c.instancesMap, key)
				continue
			}

			vCPUs := aws.ToInt32(instance.CpuOptions.CoreCount) * aws.ToInt32(instance.CpuOptions.ThreadsPerCore)
			c.instancesMap[key] = &v1.Instance{
				ID:       id,
				Name:     getInstanceTag(instance.Tags, "Name"),
				Provider: provider,
				Service:  ec2Service,
				Region:   region,
				Kind:     string(instance.InstanceType),
				Labels: v1.Labels{
					"Name":      getInstanceTag(instance.Tags, "Name"),
					"Lifecycle": string(instance.InstanceLifecycle),
					"VCPUCount": fmt.Sprint(vCPUs),
				},
			}
		}
	}
}

func getInstanceTag(tags []types.Tag, key string) string {
	for _, tag := range tags {
		if aws.ToString(tag.Key) == key {
			return aws.ToString(tag.Value)
		}
	}
	return ""
}

func buildListPaginationRequest(nextToken *string) *ec2.DescribeInstancesInput {
	return &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running", "pending"},
			},
		},
		MaxResults: aws.Int32(50),
		NextToken:  nextToken,
	}
}
