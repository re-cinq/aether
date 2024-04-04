package amazon

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/re-cinq/aether/pkg/log"
	"github.com/re-cinq/aether/pkg/providers/util"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

// Get the resource consumption of an ec2 instance
func (c *Client) GetEC2Metrics(ctx context.Context, region string, interval time.Duration) error {
	end := time.Now().UTC()
	start := end.Add(-interval)

	// Get the cpu consumption for all the instances in the region
	err := c.getEC2CPU(ctx, region, start, end, interval)
	if err != nil {
		return err
	}

	return nil
}

// Get the CPU resource consumption of an ec2 instance
func (c *Client) getEC2CPU(
	ctx context.Context,
	region string,
	start, end time.Time,
	interval time.Duration,
) error {
	// Override the region
	withRegion := func(o *cloudwatch.Options) {
		o.Region = region
	}

	logger := log.FromContext(ctx)

	period := int32(interval.Seconds())
	// validate the casting from float64 to int32
	if float64(period) != interval.Seconds() {
		return fmt.Errorf("error casting %+v to int32", interval.Seconds())
	}

	// Make the call to get the CPU metrics
	output, err := c.cloudwatch.GetMetricData(ctx, &cloudwatch.GetMetricDataInput{
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
		return err
	}

	// Loop through the result and build the intermediate awsMetric model
	for _, metric := range output.MetricDataResults {
		instanceID := aws.ToString(metric.Label)
		if instanceID == "Other" {
			return errors.New("error bad query passed to GetMetricData - instanceID not found in label")
		}

		if len(metric.Values) > 0 {
			m := v1.NewMetric(v1.CPU.String())
			m.Unit = v1.VCPU
			m.Usage = metric.Values[0]
			m.ResourceType = v1.CPU
			m.Labels = v1.Labels{
				"instanceID": instanceID,
			}

			// Update cached instance with metric
			key := util.Key(region, ec2Service, instanceID)
			instance, ok := c.instancesMap[key]
			if !ok {
				logger.Warn("instance not found in cache", "error", err, "key", key)
				continue
			}

			// ParseFloat returns 0 on failure, since that's the default
			// value of an unassigned int, store it regardless of the
			// error. This value for vCPUs is a fallback to that provided
			// by the dataset.
			if vCPUs, exists := instance.Labels["VCPUCount"]; exists {
				m.UnitAmount, err = strconv.ParseFloat(vCPUs, 64)
				if err != nil {
					logger.Error("failed to parse GCP total VCPUs", "error", err)
				}
			}
			instance.Metrics.Upsert(m)
		}
	}

	return nil
}
