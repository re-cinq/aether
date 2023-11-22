package gcp

import (
	"context"
	"fmt"
	"strconv"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	monitoringpb "cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	"google.golang.org/api/iterator"
)

var (
	/*
	* An MQL query that will return data from Google Cloud with the
	* - Instance Name
	* - Region
	* - Machine Type
	* - Reserved CPUs
	* - Utilization
	 */
	CPUQuery = `
  fetch gce_instance
  | { metric 'compute.googleapis.com/instance/cpu/utilization'
    ; metric 'compute.googleapis.com/instance/cpu/reserved_cores' }
  | outer_join 0
	| filter project_id = '%s' 
  | group_by [
        metric.instance_name, 
        metadata.system.region,
				metadata.system.machine_type,
        reserved_cores: format(t_1.value.reserved_cores, '%%f')
  ], [max(t_0.value.utilization)]
  | window %s
  | within %s
	`
)

// GCP is the structure used as the provider for Google Cloud Platform
type GCP struct {
	client    *monitoring.QueryClient
	projectID string
}

type options func(*GCP)

// WithProjectID is an option used to pass a projectID to the provider
func WithProjectID(projectID string) options {
	return func(g *GCP) {
		g.projectID = projectID
	}
}

// New returns a new instance of the GCP provider as well as a function to
// cleanup connections once done
func New(
	ctx context.Context,
	opts ...options,
) (g *GCP, teardown func(), err error) {
	// set any defaults here
	g = &GCP{}

	// overwrite any options
	for _, opt := range opts {
		opt(g)
	}

	// This allows overwriting the default client
	// google by default trys to authenticate when initilizing
	// a client therefore if we put this before running the options
	// it would try authenticate against google regardless of overwriting the
	// client
	if g.client == nil {
		c, err := monitoring.NewQueryClient(ctx)
		if err != nil {
			return nil, func() {}, err
		}
		g.client = c
	}

	// teardown is used to close relevant connections
	// and cleanup
	teardown = func() {
		g.client.Close()
	}

	return g, teardown, nil
}

// GetCPUUtilization returns the utilization for instances and is a wrapper
func (g *GCP) GetCPUForInstances(
	ctx context.Context,
	window string,
) ([]*v1.Metric, error) {
	return g.instanceMetrics(
		ctx, fmt.Sprintf(CPUQuery, g.projectID, window, window),
	)
}

// instanceMetrics runs a query on googe cloud monitoring using MQL
// and responds with a list of metrics
func (g *GCP) instanceMetrics(
	ctx context.Context,
	query string,
) ([]*v1.Metric, error) {
	var metrics []*v1.Metric
	it := g.client.QueryTimeSeries(ctx, &monitoringpb.QueryTimeSeriesRequest{
		Name:  fmt.Sprintf("projects/%s", g.projectID),
		Query: query,
	})

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		m := v1.NewMetric(resp.GetLabelValues()[0].GetStringValue())

		m.SetType(v1.CPU).SetUsagePercentage(resp.GetPointData()[0].GetValues()[0].GetDoubleValue() * 100)

		f, err := strconv.ParseFloat(resp.GetLabelValues()[3].GetStringValue(), 64)
		//TODO: we should not fail here but collect errors
		if err != nil {
			return nil, err
		}

		m.SetTotal(f)
		m.SetLabels(v1.Labels{
			"machine_type": resp.GetLabelValues()[2].GetStringValue(),
			"region":       resp.GetLabelValues()[1].GetStringValue(),
		})
		metrics = append(metrics, m)
	}
	return metrics, nil
}
