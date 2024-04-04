package amazon

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/re-cinq/aether/pkg/providers/util"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
	"github.com/stretchr/testify/assert"
)

func TestGetMetricsData(t *testing.T) {
	ctx := context.TODO()
	stubber := testtools.NewStubber()

	c := &Client{
		instancesMap: make(map[string]*v1.Instance),
		cloudwatch:   cloudwatch.NewFromConfig(*stubber.SdkConfig, func(o *cloudwatch.Options) {}),
	}

	instanceID := "i-00123456789"
	region := "test region"
	interval := 5 * time.Minute
	start := time.Date(2024, 01, 15, 20, 34, 58, 651387237, time.UTC)
	end := start.Add(interval)

	t.Run("get passing metrics data", func(t *testing.T) {
		stubber.Add(testtools.Stub{
			OperationName: "GetMetricData",
			Input: &cloudwatch.GetMetricDataInput{
				StartTime: &start,
				EndTime:   &end,
				MetricDataQueries: []types.MetricDataQuery{
					{
						Id:         aws.String(v1.CPU.String()),
						Expression: aws.String(`SELECT AVG(CPUUtilization) FROM "AWS/EC2" GROUP BY InstanceId`),
						Period:     aws.Int32(300), // 5 minutes
					},
				},
			},
			Output: &cloudwatch.GetMetricDataOutput{
				MetricDataResults: []types.MetricDataResult{
					{
						Id:     aws.String("testID"),
						Label:  aws.String("i-00123456789"),
						Values: []float64{.0000123},
					},
				},
			},
		})

		// Update the cache with a test instance to check that
		// the metrics are added
		testKey := util.Key(region, ec2Service, instanceID)
		instance := &v1.Instance{}

		// set the instance in the cache for adding a
		// metric to
		c.instancesMap[testKey] = instance

		// get the metrics and update the instance in the cache
		err := c.getEC2CPU(ctx, region, start, end, interval)
		assert.Nil(t, err)

		res, ok := c.instancesMap[testKey]
		assert.True(t, ok)
		r := res.Metrics["cpu"]

		testtools.ExitTest(stubber, t)

		expRes := v1.Metric{
			Name:         "cpu",
			Usage:        0.0000123,
			Unit:         v1.VCPU,
			ResourceType: v1.CPU,
			Labels: v1.Labels{
				"instanceID": "i-00123456789",
			},
		}

		// Comparing the metric structs fails because the time.Time timeStamp field will
		// not be equal, and since the fields are not exported it cannot be modified or
		// excluded in the compare.
		// So instead compare each field value individually

		// v1.Metric.String() creates a string of type, name, amount, and usage
		assert.Equalf(t, expRes.String(), r.String(), "Result should be: %v, got: %v", expRes, r)
		// compare labels
		assert.Equalf(t, expRes.Labels, r.Labels, "Result should be: %v, got: %v", expRes, r)
		// compare Resource Unit
		assert.Equalf(t, expRes.Unit, r.Unit, "Result should be: %v, got: %v", expRes, r)
		// emissions should not yet be calculated at this point
		assert.Equal(t, r.Emissions, v1.ResourceEmissions{})
	})

	t.Run("error getting metrics", func(t *testing.T) {
		stubber.Add(testtools.Stub{
			OperationName: "GetMetricData",
			Error:         &testtools.StubError{Err: errors.New("Testing the error is handled")},
		})

		err := c.getEC2CPU(ctx, region, start, end, interval)
		testtools.ExitTest(stubber, t)

		assert.Error(t, err)
	})
}
