package amazon

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	"k8s.io/klog/v2"
)

// Helper service to get CloudWatch data
type cloudWatchClient struct {
	client *cloudwatch.Client
}

// New cloudwatch client instance
func NewCloudWatchClient(cfg *aws.Config) *cloudWatchClient {
	emptyOptions := func(o *cloudwatch.Options) {}

	// Init the Cloudwatch client
	client := cloudwatch.NewFromConfig(*cfg, emptyOptions)

	// Make sure the initialisation was successful
	if client == nil {
		klog.Fatal("failed to create AWS CloudWatch client")
		return nil
	}

	// Return the cloudwatch service client
	return &cloudWatchClient{
		client: client,
	}
}

// Get the resource consumption of an ec2 instance
func (e *cloudWatchClient) GetEC2Metrics(region awsRegion, interval time.Duration, cache *awsCache) (map[string]v1.Instance, error) {
	serviceMetrics := make(map[string]v1.Instance)

	end := time.Now().UTC()
	start := end.Add(-interval)

	// Get the cpu consumption for all the instances in the region
	cpuMetrics, err := e.getEC2CPU(region, start, end, interval)
	if err != nil || len(cpuMetrics) < 1 {
		return serviceMetrics, err
	}
	for _, cpuMetric := range cpuMetrics {
		// load the instance metadata from the cache, because the query does not give us
		instanceMetadata := cache.Get(region, ec2Service, cpuMetric.instanceID)
		if instanceMetadata == nil {
			klog.Warningf("instance id %s is not present in the metadata, temporarily skipping collecting metrics", cpuMetric.instanceID)
			continue
		}

		// if we got here it means that we do have the instance metadata
		instanceService, exists := serviceMetrics[cpuMetric.instanceID]
		if !exists {
			// Then create a new one
			s := v1.NewInstance(cpuMetric.instanceID, provider)
			s.SetService("EC2")
			s.SetKind(instanceMetadata.kind).SetRegion(region)
			s.AddLabel("Name", instanceMetadata.name)
			serviceMetrics[cpuMetric.instanceID] = *s

			// Makes it easier to use it
			instanceService = serviceMetrics[cpuMetric.instanceID]
			instanceService.SetKind(instanceMetadata.kind).SetRegion(region)
		}

		// Build the resource
		cpu := v1.NewMetric(v1.CPU.String()).SetResourceUnit(cpuMetric.unit).SetUnitAmount(float64(instanceMetadata.coreCount))
		cpu.SetUsage(cpuMetric.value).SetType(cpuMetric.kind)

		// Update the CPU information now
		instanceService.Metrics().Upsert(cpu)
	}

	// Return the collected metrics
	return serviceMetrics, nil
}

// Get the CPU resource consumption of an ec2 instance
func (e *cloudWatchClient) getEC2CPU(region awsRegion, start, end time.Time, interval time.Duration) ([]awsMetric, error) {
	// Override the region
	withRegion := func(o *cloudwatch.Options) {
		o.Region = region
	}

	period := int32(interval.Seconds())
	// validate the casting from float64 to int32
	if float64(period) != interval.Seconds() {
		return nil, fmt.Errorf("error casting %+v to int32", interval.Seconds())
	}

	// Make the call to get the CPU metrics
	output, err := e.client.GetMetricData(context.TODO(), &cloudwatch.GetMetricDataInput{
		StartTime: &start,
		EndTime:   &end,
		MetricDataQueries: []types.MetricDataQuery{
			{
				Id:         aws.String(v1.CPU.String()),
				Expression: aws.String(`SELECT AVG(CPUUtilization) FROM "AWS/EC2" GROUP BY InstanceId`),
				Period:     aws.Int32(period),
			},
		},
	}, withRegion)
	if err != nil {
		return nil, err
	}

	// Collector
	var cpuMetrics []awsMetric

	// Loop through the result and build the intermediate awsMetric model
	for _, metric := range output.MetricDataResults {
		if len(metric.Values) > 0 {
			cpuMetric := awsMetric{
				value:      metric.Values[0],
				instanceID: aws.ToString(metric.Label),
				kind:       v1.CPU,
				unit:       v1.Core,
				name:       v1.CPU.String(),
			}

			if cpuMetric.instanceID == "Other" {
				return nil, errors.New("error bad query passed to GetMetricData - instanceID not found in label")
			}

			cpuMetrics = append(cpuMetrics, cpuMetric)
		}
	}

	return cpuMetrics, nil
}
