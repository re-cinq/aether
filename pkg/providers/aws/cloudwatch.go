package amazon

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	"k8s.io/klog/v2"
)

// Helper service to get CloudWatch data
type cloudWatchClient struct {
	client *cloudwatch.Client
}

// New cloudwatch client instance
func NewCloudWatchClient(cfg aws.Config) *cloudWatchClient {

	emptyOptions := func(o *cloudwatch.Options) {}

	// Init the Cloudwatch client
	client := cloudwatch.NewFromConfig(cfg, emptyOptions)

	// Make sure the initialisation was successful
	if client == nil {
		klog.Fatal("failed to create AWS CloudWatch client")
		return nil
	}

	// Return the ec2 service
	return &cloudWatchClient{
		client: client,
	}

}

// Get the resource consumption of an ec2 instance
func (e *cloudWatchClient) getEc2Metrics(region awsRegion, instanceId string) []v1.Service {

	var serviceMetrics []v1.Service

	// Build the service
	serviceMetric := v1.NewService(instanceId, awsProvider).SetRegion(region)

	// Get the cpu consumption
	if cpuMetric := e.getEc2Cpu(region, instanceId); cpuMetric != nil {
		serviceMetric.Metrics().Upsert(cpuMetric)
		serviceMetrics = append(serviceMetrics, serviceMetric.Build())
	}

	// Return the collected metrics
	return serviceMetrics

}

// Get the CPU resource consumption of an ec2 instance
func (e *cloudWatchClient) getEc2Cpu(region awsRegion, instanceId string) *v1.Metric {

	return nil
}
