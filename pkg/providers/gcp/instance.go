package gcp

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	monitoringpb "cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	"google.golang.org/api/iterator"
)

var (
	/*
	* An MQL query that will return data from Google Cloud with the
	* - Instance Name
	* - Region
	* - Zone
	* - Machine Type
	* - Reserved CPUs
	* - Utilization
	* NOTE: Using reserved CPUs as vCPUs, because they are equivalent for visible
	* vCPUs within a guest instance, except for shared-core machines:
	* https://cloud.google.com/monitoring/api/metrics_gcp
	 */
	CPUQuery = `
  fetch gce_instance
  | { metric 'compute.googleapis.com/instance/cpu/utilization'
    ; metric 'compute.googleapis.com/instance/cpu/reserved_cores' }
  | outer_join 0
	| filter project_id = '%s' 
  | group_by [		
    resource.instance_id,
  	metric.instance_name,
		metadata.system.region,
		resource.zone,
		metadata.system.machine_type,
    reserved_cores: format(t_1.value.reserved_cores, '%%f')
  ], [max(t_0.value.utilization)]
  | window %s
  | within %s
	`
	/*
	* An MQL query that will return memory data from Google Cloud with the
	* - Instance Name
	* - Region
	* - Zone
	* - Machine Type
	* - Memory Usage
	* NOTE: According to Google the 'ram_used' metric is only available for
	* e2-xxxx instances, which means that we can get memory usage for other types
	* of VM's
	 */
	MEMQuery = `
	fetch gce_instance
	| metric 'compute.googleapis.com/instance/memory/balloon/ram_used'
	| filter project_id = '%s'
	| group_by [
	  resource.instance_id,
	  metric.instance_name,
		metadata.system.region,
		resource.zone,
		metadata.system.machine_type,
	], [max(value.ram_used)]
	| window %s
	| within %s
	`
)

// instanceMetrics runs a query on googe cloud monitoring using MQL
// and responds with a list of metrics
func (g *GCP) instanceMemoryMetrics(
	ctx context.Context,
	project, query string,
) ([]*v1.Metric, error) {
	var metrics []*v1.Metric

	it := g.monitoring.QueryTimeSeries(ctx, &monitoringpb.QueryTimeSeriesRequest{
		Name:  fmt.Sprintf("projects/%s", project),
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

		// This is dependant on the MQL query
		// label ordering
		instanceID := resp.GetLabelValues()[0].GetStringValue()
		instanceName := resp.GetLabelValues()[1].GetStringValue()
		region := resp.GetLabelValues()[2].GetStringValue()
		zone := resp.GetLabelValues()[3].GetStringValue()
		instanceType := resp.GetLabelValues()[4].GetStringValue()

		m := v1.NewMetric("memory")
		m.SetResourceUnit(v1.Gb)
		m.SetType(v1.Memory).SetUsage(
			// convert bytes to GBs
			float64(resp.GetPointData()[0].GetValues()[0].GetInt64Value()) / 1024 / 1024 / 1024,
		)

		// TODO: we should not fail here but collect errors
		if err != nil {
			slog.Error("failed to parse GCP metric", "error", err)
			continue
		}

		m.SetLabels(v1.Labels{
			"id":           instanceID,
			"name":         instanceName,
			"region":       region,
			"zone":         zone,
			"machine_type": instanceType,
		})
		metrics = append(metrics, m)
	}
	return metrics, nil
}

// instanceCPUMetrics runs a query on googe cloud monitoring using MQL
// and responds with a list of CPU metrics
func (g *GCP) instanceCPUMetrics(
	ctx context.Context,
	project, query string,
) ([]*v1.Metric, error) {
	var metrics []*v1.Metric

	it := g.monitoring.QueryTimeSeries(ctx, &monitoringpb.QueryTimeSeriesRequest{
		Name:  fmt.Sprintf("projects/%s", project),
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

		// This is dependant on the MQL query
		// label ordering
		instanceID := resp.GetLabelValues()[0].GetStringValue()
		instanceName := resp.GetLabelValues()[1].GetStringValue()
		region := resp.GetLabelValues()[2].GetStringValue()
		zone := resp.GetLabelValues()[3].GetStringValue()
		instanceType := resp.GetLabelValues()[4].GetStringValue()
		totalVCPUs := resp.GetLabelValues()[5].GetStringValue()

		m := v1.NewMetric("cpu")
		m.SetResourceUnit(v1.VCPU)
		m.SetType(v1.CPU).SetUsage(
			// translate fraction to a percentage
			resp.GetPointData()[0].GetValues()[0].GetDoubleValue() * 100,
		)

		f, err := strconv.ParseFloat(totalVCPUs, 64)
		// TODO: we should not fail here but collect errors
		if err != nil {
			slog.Error("failed to parse GCP metric", "error", err)
			continue
		}

		m.SetUnitAmount(f)
		m.SetLabels(v1.Labels{
			"id":           instanceID,
			"name":         instanceName,
			"region":       region,
			"zone":         zone,
			"machine_type": instanceType,
		})
		metrics = append(metrics, m)
	}
	return metrics, nil
}
