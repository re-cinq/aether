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

	input := &cloudwatch.GetMetricDataInput{
		StartTime: &start,
		EndTime:   &end,
		MetricDataQueries: []types.MetricDataQuery{
			{
				Id:     aws.String(v1.CPU.String()),
				Period: aws.Int32(300), // 5 minutes
			},
		},
	}

	// first stub call to get cpu metrics
	input.MetricDataQueries[0].Expression = aws.String(cpuExpression)
	stubber.Add(testtools.Stub{
		OperationName: "GetMetricData",
		Input:         input,
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

	// second stub call to get memory metrics
	input.MetricDataQueries[0].Expression = aws.String(`SELECT AVG(mem_used_percent) FROM SCHEMA(CWAgent, InstanceId) GROUP BY InstanceId`)
	stubber.Add(testtools.Stub{
		OperationName: "GetMetricData",
		Input:         input,
		Output: &cloudwatch.GetMetricDataOutput{
			MetricDataResults: []types.MetricDataResult{
				{
					Id:     aws.String("testID"),
					Label:  aws.String("i-00123456789"),
					Values: []float64{.27},
				},
			},
		},
	})

	// third stub erroring the GetMetricData call for both cpu and memory
	stubber.Add(testtools.Stub{
		OperationName: "GetMetricData",
		Error:         &testtools.StubError{Err: errors.New("Testing the error is handled")},
	})

	t.Run("pass CPU metrics data", func(t *testing.T) {
		// Update the cache with a test instance to check that
		// the metrics are added
		testKey := util.Key(region, ec2Service, instanceID)
		instance := &v1.Instance{}

		// set the instance in the cache for adding a
		// metric to
		c.instancesMap[testKey] = instance

		// get the metrics and update the instance in the cache
		err := c.cpuMetrics(ctx, region, input)
		assert.Nil(t, err)

		res, ok := c.instancesMap[testKey]
		assert.True(t, ok)
		r := res.Metrics["cpu"]

		expRes := v1.Metric{
			Name:         "cpu",
			Usage:        0.0000123,
			Unit:         v1.VCPU,
			ResourceType: v1.CPU,
			Labels: v1.Labels{
				"instanceID": "i-00123456789",
			},
		}

		// v1.Metric.String() creates a string of type, name, amount, and usage
		assert.Equalf(t, expRes.String(), r.String(), "Result should be: %v, got: %v", expRes, r)
		// compare labels
		assert.Equalf(t, expRes.Labels, r.Labels, "Result should be: %v, got: %v", expRes, r)
		// compare Resource Unit
		assert.Equalf(t, expRes.Unit, r.Unit, "Result should be: %v, got: %v", expRes, r)
		// emissions should not yet be calculated at this point
		assert.Equal(t, r.Emissions, v1.ResourceEmissions{})
	})

	t.Run("pass memory metrics data", func(t *testing.T) {
		// Update the cache with a test instance to check that
		// the metrics are added
		testKey := util.Key(region, ec2Service, instanceID)
		instance := &v1.Instance{}

		// set the instance in the cache for adding a
		// metric to
		c.instancesMap[testKey] = instance

		// get the metrics and update the instance in the cache
		err := c.memoryMetrics(ctx, region, input)
		assert.Nil(t, err)

		res, ok := c.instancesMap[testKey]
		assert.True(t, ok)
		r := res.Metrics["memory"]

		expRes := v1.Metric{
			Name:         "memory",
			Usage:        0.27,
			Unit:         v1.GB,
			ResourceType: v1.Memory,
			Labels: v1.Labels{
				"instanceID": "i-00123456789",
				"region":     region,
				"name":       "i-00123456789",
			},
		}

		// v1.Metric.String() creates a string of type, name, amount, and usage
		assert.Equalf(t, expRes.String(), r.String(), "Result should be: %v, got: %v", expRes, r)
		// compare labels
		assert.Equalf(t, expRes.Labels, r.Labels, "Result should be: %v, got: %v", expRes, r)
		// compare Resource Unit
		assert.Equalf(t, expRes.Unit, r.Unit, "Result should be: %v, got: %v", expRes, r)
		// emissions should not yet be calculated at this point
		assert.Equal(t, r.Emissions, v1.ResourceEmissions{})
	})

	t.Run("error getting memory metrics", func(t *testing.T) {
		err := c.memoryMetrics(ctx, region, input)
		assert.Error(t, err)
	})

	t.Run("error getting CPU metrics", func(t *testing.T) {
		err := c.cpuMetrics(ctx, region, input)
		assert.Error(t, err)
	})

	t.Run("test all stubs called", func(t *testing.T) {
		err := stubber.VerifyAllStubsCalled()
		assert.Nil(t, err)
	})
}

// silly little test to be sure if the query changes it's
// intentional
func TestQueriesDontChange(t *testing.T) {
	cpu := `SELECT AVG(CPUUtilization) FROM "AWS/EC2" GROUP BY InstanceId`
	assert.Equal(t, cpu, cpuExpression)

	mem := `SELECT AVG(mem_used_percent) FROM SCHEMA(CWAgent, InstanceId) GROUP BY InstanceId`
	assert.Equal(t, mem, memExpression)
}
