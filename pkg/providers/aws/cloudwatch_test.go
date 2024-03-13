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
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	"github.com/stretchr/testify/assert"
)

func TestGetMetricsData(t *testing.T) {
	stubber := testtools.NewStubber()
	client := NewCloudWatchClient(context.TODO(), stubber.SdkConfig)

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

		res, err := client.getEC2CPU(region, start, end, interval)
		testtools.ExitTest(stubber, t)

		expRes := v1.NewMetric("cpu").SetUsage(float64(.0000123)).SetResourceUnit("vCPU").SetType(v1.CPU)
		expRes.SetLabels(map[string]string{
			"instanceID": "i-00123456789",
		})

		// Comparing the metric structs fails because the time.Time timeStamp field will
		// not be equal, and since the fields are not exported it cannot be modified or
		// excluded in the compare.
		// So instead compare each field value individually

		// v1.Metric.String() creates a string of type, name, amount, and usage
		assert.Equalf(t, expRes.String(), res[0].String(), "Result should be: %v, got: %v", expRes, res)
		// compare labels
		assert.Equalf(t, expRes.Labels(), res[0].Labels(), "Result should be: %v, got: %v", expRes, res)
		// compare Resource Unit
		assert.Equalf(t, expRes.Unit(), res[0].Unit(), "Result should be: %v, got: %v", expRes, res)
		// emissions should not yet be calculated at this point
		assert.Equal(t, res[0].Emissions(), &v1.ResourceEmissions{})
		// check no error
		assert.Nil(t, err)
	})

	t.Run("error getting metrics", func(t *testing.T) {
		stubber.Add(testtools.Stub{
			OperationName: "GetMetricData",
			Error:         &testtools.StubError{Err: errors.New("Testing the error is handled")},
		})

		res, err := client.getEC2CPU(region, start, end, interval)
		testtools.ExitTest(stubber, t)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
}
