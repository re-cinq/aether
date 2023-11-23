// Contains a set of method for getting EC2 information
package amazon

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"k8s.io/klog/v2"
)

// Helper service to get EC2 data
type ec2Client struct {
	client *ec2.Client
	cache  *awsCache
}

// New instance
func NewEc2Client(cfg aws.Config) *ec2Client {

	emptyOptions := func(o *ec2.Options) {}

	// Init the EC2 client
	client := ec2.NewFromConfig(cfg, emptyOptions)

	// Make sure the initialisation was successful
	if client == nil {
		klog.Fatal("failed to create AWS EC2 client")
		return nil
	}

	// Return the ec2 service
	return &ec2Client{
		client: client,
		cache:  newAWSCache(),
	}

}

func (e *ec2Client) Cache() *awsCache {
	return e.cache
}

// refresh stores all the instances for a specific region in cache
func (e *ec2Client) Refresh(region awsRegion) {

	// Override the region
	withRegion := func(o *ec2.Options) {
		o.Region = region
	}

	// First request
	output, err := e.client.DescribeInstances(context.TODO(), buildListPaginationRequest(nil), withRegion)
	if err != nil {
		klog.Errorf("failed to retrieve ec2 instances %s", err)
	}

	if output == nil {
		klog.Errorf("failed to retrieve the list of EC2 instances from %s", region)
		return
	}

	// Make sure the first request is successful

	// Collect all the responses for all the pages
	instances := []ec2.DescribeInstancesOutput{*output}

	for output.NextToken != nil {
		output, err = e.client.DescribeInstances(context.TODO(), buildListPaginationRequest(output.NextToken), withRegion)
		if err != nil {
			klog.Errorf("failed to retrieve ec2 instances %s", err)
		}
		// Collect the response
		if output != nil {
			instances = append(instances, *output)
		}
	}

	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {

			e.cache.Add(newAWSResource(
				region,
				ec2Service,
				*instance.InstanceId,
				string(instance.InstanceType),
				string(instance.InstanceLifecycle),
				getInstanceTag(instance.Tags, "Name"),
				int(*instance.CpuOptions.CoreCount),
			))
		}
	}

}

func getInstanceTag(tags []types.Tag, key string) string {
	for _, tag := range tags {
		if *tag.Key == key {
			return *tag.Value
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
