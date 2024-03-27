package gcp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/re-cinq/aether/pkg/config"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func withMonitoringTestClient(mc *monitoring.QueryClient) options {
	return func(c *Client) {
		c.monitoring = mc
	}
}

func withInstancesTestClient(ic *compute.InstancesClient) options {
	return func(c *Client) {
		c.instances = ic
	}
}

type fakeMonitoringServer struct {
	monitoringpb.UnimplementedQueryServiceServer
	// Response that will return from the fake server
	Response *monitoringpb.QueryTimeSeriesResponse
	// If error is set the server will return an error
	Error error
}

type fakeInstancesServer struct {
	computepb.UnimplementedInstancesServer
}

// setupFakeServer is used to setup a fake GRPC server to hanlde test requests
func setupFakeServer(
	m *fakeMonitoringServer,
	i *fakeInstancesServer,
) (*string, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}
	gsrv := grpc.NewServer()
	monitoringpb.RegisterQueryServiceServer(gsrv, m)
	computepb.RegisterInstancesServer(gsrv, i)
	fakeServerAddr := l.Addr().String()
	go func() {
		if err := gsrv.Serve(l); err != nil {
			panic(err)
		}
	}()
	return &fakeServerAddr, nil
}

func (f *fakeMonitoringServer) QueryTimeSeries(
	ctx context.Context,
	req *monitoringpb.QueryTimeSeriesRequest,
) (*monitoringpb.QueryTimeSeriesResponse, error) {
	if f.Error != nil {
		return nil, f.Error
	}
	return f.Response, nil
}

type testMetric struct {
	Name       string
	UnitAmount float64
	Type       v1.ResourceType
	Labels     v1.Labels
	Usage      float64
}

var defaultLabelValues = []*monitoringpb.LabelValue{
	{
		Value: &monitoringpb.LabelValue_StringValue{StringValue: "my-instance-id"},
	},
	{
		Value: &monitoringpb.LabelValue_StringValue{StringValue: "foobar"},
	},
	{
		Value: &monitoringpb.LabelValue_StringValue{StringValue: "europe-west-1"},
	},
	{
		Value: &monitoringpb.LabelValue_StringValue{StringValue: "europe-west"},
	},
	{
		Value: &monitoringpb.LabelValue_StringValue{StringValue: "e2-medium"},
	},
	{
		Value: &monitoringpb.LabelValue_StringValue{StringValue: "2.000000"},
	},
}

type TestScenario struct {
	description         string
	scenariotype        string
	responsePointData   []*monitoringpb.TimeSeriesData_PointData
	responseLabelValues []*monitoringpb.LabelValue
	err                 error
	expectedResponse    []*testMetric
	query               string
}

// RunTestData is a helper function to run test scenarios
func RunTestData(t *testing.T, testdata []TestScenario) {
	t.Helper()
	assert := require.New(t)
	ctx := context.TODO()

	for i := range testdata {
		t.Run(testdata[i].description, func(t *testing.T) {
			testResp := &monitoringpb.QueryTimeSeriesResponse{
				TimeSeriesData: []*monitoringpb.TimeSeriesData{
					{
						LabelValues: defaultLabelValues,
						PointData:   testdata[i].responsePointData,
					},
				},
			}
			if len(testdata[i].responseLabelValues) != 0 {
				testResp.TimeSeriesData[0].LabelValues = testdata[i].responseLabelValues
			}

			fakeMonitoringServer := &fakeMonitoringServer{
				Response: testResp,
				Error:    testdata[i].err,
			}

			fakeInstancesServer := &fakeInstancesServer{}

			addr, err := setupFakeServer(fakeMonitoringServer, fakeInstancesServer)
			assert.NoError(err)

			m, err := monitoring.NewQueryClient(ctx,
				option.WithEndpoint(*addr),
				option.WithoutAuthentication(),
				option.WithGRPCDialOption(grpc.WithTransportCredentials(
					insecure.NewCredentials(),
				)),
			)
			assert.NoError(err)

			in, err := compute.NewInstancesRESTClient(ctx,
				option.WithEndpoint(*addr),
				option.WithoutAuthentication(),
				option.WithGRPCDialOption(grpc.WithTransportCredentials(
					insecure.NewCredentials(),
				)),
			)
			assert.NoError(err)

			g, teardown, err := New(ctx,
				&config.Account{},
				withMonitoringTestClient(m),
				withInstancesTestClient(in),
			)
			assert.NoError(err)
			defer teardown()

			var resp []*v1.Metric

			switch testdata[i].scenariotype {
			case "cpu":
				resp, err = g.instanceCPUMetrics(ctx, "", testdata[i].query)
			case "memory":
				resp, err = g.instanceMemoryMetrics(ctx, "", testdata[i].query)
			}

			if testdata[i].err == nil {
				assert.NoError(err)
				for i, r := range resp {
					assert.Equal(testdata[i].expectedResponse[i].Labels, r.Labels)
					assert.Equal(testdata[i].expectedResponse[i].Type, r.ResourceType)
					assert.Equal(testdata[i].expectedResponse[i].Usage, r.Usage)
					assert.Equal(testdata[i].expectedResponse[i].UnitAmount, r.UnitAmount)
				}
			} else {
				assert.Equal(
					fmt.Sprintf("%s", err),
					fmt.Sprintf("rpc error: code = Unknown desc = %s", testdata[i].err),
				)
			}
		})
	}
}

func TestGetCPUMetrics(t *testing.T) {
	st := "cpu"
	testdata := []TestScenario{
		{
			description:  "cpu metrics",
			scenariotype: st,
			query:        fmt.Sprintf(CPUQuery, "foobar", "5m", "5m"),
			responsePointData: []*monitoringpb.TimeSeriesData_PointData{
				{
					Values: []*monitoringpb.TypedValue{
						{
							Value: &monitoringpb.TypedValue_DoubleValue{
								DoubleValue: 0.01,
							},
						},
					},
				},
			},
			expectedResponse: []*testMetric{
				{
					Type: v1.CPU,
					Labels: v1.Labels{
						"id":           "my-instance-id",
						"machine_type": "e2-medium",
						"name":         "foobar",
						"region":       "europe-west-1",
						"zone":         "europe-west",
					},
					Usage:      1,
					UnitAmount: 2.0000,
				},
			},
		},
		{
			description:  "error occurs in query",
			scenariotype: st,
			query:        fmt.Sprintf(CPUQuery, "foobar", "5m", "5m"),
			err:          errors.New("random error occurred cpu query"),
		},
	}
	RunTestData(t, testdata)
}

func TestInstanceMemoryMetrics(t *testing.T) {
	st := "memory"
	testdata := []TestScenario{
		{
			description:  "memory metrics returned",
			scenariotype: st,
			query:        fmt.Sprintf(MEMQuery, "foobar", "5m", "5m"),
			responsePointData: []*monitoringpb.TimeSeriesData_PointData{
				{
					Values: []*monitoringpb.TypedValue{
						{
							Value: &monitoringpb.TypedValue_Int64Value{
								// 10GB
								Int64Value: 10 * 1024 * 1024 * 1024,
							},
						},
					},
				},
			},
			expectedResponse: []*testMetric{
				{
					Type: v1.Memory,
					Labels: v1.Labels{
						"id":           "my-instance-id",
						"machine_type": "e2-medium",
						"name":         "foobar",
						"region":       "europe-west-1",
						"zone":         "europe-west",
					},
					Usage:      10.0,
					UnitAmount: 0.0000,
				},
			},
		},
		{
			description:  "error occurs in query",
			scenariotype: st,
			query:        fmt.Sprintf(MEMQuery, "foobar", "5m", "5m"),
			err:          errors.New("random error occurred in memory query"),
		},
	}
	RunTestData(t, testdata)
}
