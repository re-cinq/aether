package amazon

import (
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
	client := NewCloudWatchClient(stubber.SdkConfig)

	start := time.Date(2024, 01, 15, 20, 34, 58, 651387237, time.UTC)
	end := start.Add(5 * time.Minute)
	region := "test region"

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
						Values: []float64{1.00},
					},
				},
			},
		})

		res, err := client.getEc2Cpu(region, start, end)
		testtools.ExitTest(stubber, t)

		expRes := []awsMetric{
			{
				name:       "cpu",
				kind:       "cpu",
				unit:       "core",
				value:      1,
				instanceID: "i-00123456789",
			},
		}

		assert.Equalf(t, expRes, res, "Result should be: %v, got: %v", expRes, res)
		assert.Nil(t, err)
	})

	t.Run("error getting metrics", func(t *testing.T) {
		stubber.Add(testtools.Stub{
			OperationName: "GetMetricData",
			Error:         &testtools.StubError{Err: errors.New("Testing the error is handled")},
		})

		res, err := client.getEc2Cpu(region, start, end)
		testtools.ExitTest(stubber, t)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
}
