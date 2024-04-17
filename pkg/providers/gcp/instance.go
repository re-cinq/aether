package gcp

import (
	"context"
	"fmt"
	"strconv"

	monitoringpb "cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/re-cinq/aether/pkg/log"
	"github.com/re-cinq/aether/pkg/providers/util"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
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
func (c *Client) memoryMetrics(ctx context.Context, project, query string) error {
	logger := log.FromContext(ctx)

	it := c.monitoring.QueryTimeSeries(ctx, &monitoringpb.QueryTimeSeriesRequest{
		Name:  fmt.Sprintf("projects/%s", project),
		Query: query,
	})

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		// This is dependant on the MQL query
		// label ordering
		instanceID := resp.GetLabelValues()[0].GetStringValue()
		instanceName := resp.GetLabelValues()[1].GetStringValue()
		region := resp.GetLabelValues()[2].GetStringValue()
		zone := resp.GetLabelValues()[3].GetStringValue()
		instanceType := resp.GetLabelValues()[4].GetStringValue()

		m := v1.NewMetric(v1.Memory.String())
		m.Unit = v1.GB
		m.ResourceType = v1.Memory
		// convert Bytes to GB
		m.UnitAmount = float64(resp.GetPointData()[0].GetValues()[0].GetInt64Value()) / 1024 / 1024 / 1024
		m.Labels = v1.Labels{
			"id":           instanceID,
			"name":         instanceName,
			"region":       region,
			"zone":         zone,
			"machine_type": instanceType,
		}

		// Get the stored instance, update the metric and restore
		// the instance in the cache
		key := util.Key(zone, service, instanceName)
		instance, ok := c.instancesMap[key]
		if !ok {
			logger.Warn("instance not found in cache", "error", err, "key", key)
			continue
		}

		instance.Metrics.Upsert(m)
	}
	return nil
}

// instanceCPUMetrics runs a query on googe cloud monitoring using MQL
// and responds with a list of CPU metrics
func (c *Client) cpuMetrics(ctx context.Context, project, query string) error {
	logger := log.FromContext(ctx)

	it := c.monitoring.QueryTimeSeries(ctx, &monitoringpb.QueryTimeSeriesRequest{
		Name:  fmt.Sprintf("projects/%s", project),
		Query: query,
	})

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		// This is dependant on the MQL query
		// label ordering
		instanceID := resp.GetLabelValues()[0].GetStringValue()
		instanceName := resp.GetLabelValues()[1].GetStringValue()
		region := resp.GetLabelValues()[2].GetStringValue()
		zone := resp.GetLabelValues()[3].GetStringValue()
		instanceType := resp.GetLabelValues()[4].GetStringValue()
		totalVCPUs := resp.GetLabelValues()[5].GetStringValue()

		m := v1.NewMetric(v1.CPU.String())
		m.Unit = v1.VCPU
		m.ResourceType = v1.CPU

		// translate fraction to a percentage
		m.Usage = resp.GetPointData()[0].GetValues()[0].GetDoubleValue() * 100

		// ParseFloat returns 0 on failure, since that's the default
		// value of an unassigned int, store it regardless of the
		// error. This value for vCPUs is a fallback to that provided
		// by the dataset.
		m.UnitAmount, err = strconv.ParseFloat(totalVCPUs, 64)
		if err != nil {
			logger.Error("failed to parse GCP total VCPUs", "error", err)
		}

		m.Labels = v1.Labels{
			"id":           instanceID,
			"name":         instanceName,
			"region":       region,
			"zone":         zone,
			"machine_type": instanceType,
		}

		// Get the cached instance, update the metric and restore
		// the instance in the cache
		key := util.Key(zone, service, instanceName)
		instance, ok := c.instancesMap[key]
		if !ok {
			logger.Warn("instance not found in cache", "error", err, "key", key)
			continue
		}

		// We are updating the instance based on the stored pointer.
		// Since it's by reference, we don't need to restore the
		// object in the map
		instance.Metrics.Upsert(m)
	}
	return nil
}
