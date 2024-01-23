package calculator

import (
	"errors"
	"fmt"
	"time"

	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	"k8s.io/klog/v2"
)

// TODO add links / sources for where numbers and calculations are gathered

// lifespan is the default amount of years a server is in use before being
// replaced in a datacenter
const lifespan = 4

type calculate struct {
	minWatts      float64
	maxWatts      float64
	chip          float64
	pue           float64
	gridCO2e      float64
	totalEmbodied float64
}

// TODO: This was borrowed from CCF calculations, do more research
const memCoefficient = 0.000392 // Kilowatt Hour / Gigabyte Hour.

// operationalEmissions determines the correct function to run to calculate the
// operational emissions for the metric type
func (c *calculate) operationalEmissions(metric *v1.Metric, interval time.Duration) (float64, error) {
	switch metric.Name() {
	case v1.CPU.String():
		return c.cpu(metric, interval)
	case v1.Memory.String():
		return c.memory(metric, interval), nil
	case v1.Storage.String():
		return float64(-1), errors.New("error storage is not yet being calculated")
	case v1.Network.String():
		return float64(-1), errors.New("error networking is not yet being calculated")
	default:
		return 0, fmt.Errorf("error metric not found to be calculated: %+v", metric)
	}
}

// cpu are the emissions released from the machines the service is
// running on based on architecture and utilization.
func (c *calculate) cpu(m *v1.Metric, interval time.Duration) (float64, error) {
	// Check that number of cores is set
	if m.UnitAmount() == 0 {
		return 0, errors.New("error Cores set to 0, this should never be the case")
	}

	// vCPUHours is the amount of cores on the machine multiplied by the interval of time
	// for 1 hour. For example, if the machine has 4 cores and the interval of time is
	// 5 minutes: The hourly time is 5/60 (0.083333333) * 4 cores = 0.333333333.
	//nolint:unconvert //conversion to minutes does affect calculation
	vCPUHours := m.UnitAmount() * (float64(interval.Minutes()) / float64(60))

	// Average Watts is the average energy consumption of the service. It is based on
	// CPU utilization and Minimum and Maximum wattage of the server. If the machine
	// architecture is unknown the Min and Max wattage is the average of all machines
	// for that provider, and is supplied in the provider defaults. This is being
	// handled in the types/factors package (the point of reading in coefficient data).
	if m.Usage() == 0 {
		// TODO: Should this be an error?
		klog.Warning("CPU metric has no usage")
	}

	avgWatts := c.minWatts + m.Usage()*(c.maxWatts-c.minWatts)

	// Operational Emissions are calculated by multiplying the avgWatts, vCPUHours, PUE,
	// and region grid CO2e. The PUE is collected from the providers. The CO2e grid data
	// is the electrical grid emissions for the region at the specified time.
	return avgWatts * vCPUHours * c.pue * c.gridCO2e, nil
}

// memory calculated the CO2e emissions from the memory of
// running services on a machine. It is calculated in GB/hours
func (c *calculate) memory(m *v1.Metric, interval time.Duration) float64 {

	// GCP: Kilowatt hours = Memory usage (GB-Hours) x Memory coefficient
	// RAM usage determined by operating voltage and clock speed (MHz)
	// AWS uses DDR4 and DDR5
	return (m.Usage() * (float64(interval.Minutes()) / float64(60))) * memCoefficient
}

// EmbodiedEmissions are the released emissions of production and destruction of the
// hardware
func (c *calculate) embodiedEmissions(interval time.Duration) float64 {
	// Total Embodied is the total emissions for a server to be produced, including
	// additional emmissions for added DRAM, CPUs, GPUS, and storage. This is divided
	// by the expected lifespan of the server to get the annual emissions.
	annualEmbodied := c.totalEmbodied / lifespan

	// The embodied emissions need to be calculated for the measurement interval, so the
	// annual emissions further divided to the interval minutes.
	//nolint:unconvert //conversion to minutes does affect calculation
	return annualEmbodied / float64(365) / float64(24) / float64(60) * float64(interval.Minutes())
}
