package gcp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/re-cinq/cloud-carbon/pkg/config"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func withTestClient(c *monitoring.QueryClient) options {
	return func(g *GCP) {
		g.client = c
	}
}

type fakeMonitoringServer struct {
	monitoringpb.UnimplementedQueryServiceServer
	// Response that will return from the fake server
	Response *monitoringpb.QueryTimeSeriesResponse
	// If error is set the server will return an error
	Error error
}

// setupFakeServer is used to setup a fake GRPC server to hanlde test requests
func setupFakeServer(f *fakeMonitoringServer) (*string, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}
	gsrv := grpc.NewServer()
	monitoringpb.RegisterQueryServiceServer(gsrv, f)
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

var defaultLabelValues = []*monitoringpb.LabelValue{
	{
		Value: &monitoringpb.LabelValue_StringValue{StringValue: "foobar"},
	},
	{
		Value: &monitoringpb.LabelValue_StringValue{StringValue: "europe-west-1"},
	},
	{
		Value: &monitoringpb.LabelValue_StringValue{StringValue: "e2-medium"},
	},
	{
		Value: &monitoringpb.LabelValue_StringValue{StringValue: "2.000000"},
	},
}

func TestGetCPUMetrics(t *testing.T) {
	assert := require.New(t)
	ctx := context.TODO()

	//TODO see if we can use v1.Metric instead of this
	type testMetric struct {
		Name   string
		Total  float64
		Type   v1.ResourceType
		Labels v1.Labels
		Usage  v1.Percentage
	}
	testdata := []struct {
		description         string
		responsePointData   []*monitoringpb.TimeSeriesData_PointData
		responseLabelValues []*monitoringpb.LabelValue
		err                 error
		expectedResponse    []*testMetric
		query               string
	}{
		{
			description: "query for count metrics",
			query:       fmt.Sprintf(CPUQuery, "foobar", "5m", "5m"),
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
						"machine_type": "e2-medium",
						"region":       "europe-west-1",
					},
					Usage: 1,
					Total: 2.0000,
				},
			},
		},
		{
			description: "error occurs in query",
			query:       fmt.Sprintf(CPUQuery, "foobar", "5m", "5m"),
			err:         errors.New("random error occurred"),
		},
	}

	for _, test := range testdata {
		t.Run(test.description, func(t *testing.T) {
			testResp := &monitoringpb.QueryTimeSeriesResponse{
				TimeSeriesData: []*monitoringpb.TimeSeriesData{
					{
						LabelValues: defaultLabelValues,
						PointData:   test.responsePointData,
					},
				},
			}
			if len(test.responseLabelValues) != 0 {
				testResp.TimeSeriesData[0].LabelValues = test.responseLabelValues
			}

			fakeMonitoringServer := &fakeMonitoringServer{
				Response: testResp,
				Error:    test.err,
			}

			addr, err := setupFakeServer(fakeMonitoringServer)
			assert.NoError(err)

			client, err := monitoring.NewQueryClient(ctx,
				option.WithEndpoint(*addr),
				option.WithoutAuthentication(),
				option.WithGRPCDialOption(grpc.WithTransportCredentials(
					insecure.NewCredentials(),
				)),
			)
			assert.NoError(err)

			g, teardown, err := New(config.Account{}, newGCPCache(), withTestClient(client))
			assert.NoError(err)
			defer teardown()

			resp, err := g.instanceMetrics(ctx, test.query)
			if test.err == nil {
				assert.NoError(err)
				for i, r := range resp {
					assert.Equal(test.expectedResponse[i].Labels, r.Labels())
					assert.Equal(test.expectedResponse[i].Type, r.Type())
					assert.Equal(test.expectedResponse[i].Usage, r.Usage())
					assert.Equal(test.expectedResponse[i].Total, r.Total())
				}
			} else {
				assert.Equal(
					fmt.Sprintf("%s", err),
					fmt.Sprintf("rpc error: code = Unknown desc = %s", test.err),
				)
			}
		})
	}
}
