package calculator

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cnkei/gospline"
	"github.com/re-cinq/aether/pkg/log"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
	data "github.com/re-cinq/emissions-data/pkg/types/v2"
)

type parameters struct {
	grid           float64
	pue            float64
	metric         *v1.Metric
	factors        *data.Instance
	embodiedFactor float64
}

// operationalEmissions determines the correct function to run to calculate the
// operational emissions for the metric type which stores the energy consumption
// and the carbon emissions in the metric
func operationalEmissions(ctx context.Context, interval time.Duration, p *parameters) error {
	switch p.metric.Name {
	case v1.CPU.String():
		return cpu(ctx, interval, p)
	case v1.Memory.String():
		return memory(ctx, p)
	case v1.Storage.String():
		return errors.New("error storage is not yet being calculated")
	case v1.Network.String():
		return errors.New("error networking is not yet being calculated")
	default:
		return fmt.Errorf("error metric not supported: %+v", p.metric.Name)
	}
}

// cpu calculates the CO2e operational emissions for the CPU utilization of
// a Cloud VM instance over an interval of time.
//
// The initial calculation uses the wattage conversion factor based on the turbostat and
// turbostress to stress test the CPU on baremetal servers as inspired by Teads.
// More information can be found in our docs/METHODOLOGIES.md
func cpu(ctx context.Context, interval time.Duration, p *parameters) error {
	logger := log.FromContext(ctx)

	vCPU := p.vCPU
	// vCPU are virtual CPUs that are mapped to physical cores (a core is a physical
	// component to the CPU the VM is running on). If vCPU from the dataset (p.vCPU)
	// is not found, get the number of vCPUs from the metric collected from the query
	if vCPU == 0 {
		if p.metric.UnitAmount == 0 {
			return errors.New("error vCPU set to 0")
		}
		vCPU = p.metric.UnitAmount
	}

	// vCPUHours represents the count of virtual CPUs within a specific time frame.
	// To get vCPUHours, we first get the interval in hours and multiply that by the
	// number of vCPUs.
	// For example, if the machine has 4 vCPUs and an interval of time of 5 minutes
	// The hourly time is 5/60 (0.083333333) * 4 vCPU = 0.33333334
	vCPUHours := (interval.Minutes() / float64(60)) * vCPU

	// energy is the CPU energy consumption in kilowatts.
	// If pkgWatt values exist from the dataset, then use cubic spline interpolation
	// to calculate the wattage based on utilization.
	usage, err := cubicSplineInterpolation(p.factors.PkgWatt, p.metric.Usage)
	if err != nil {
		return err
	}

	p.metric.Energy = usage * vCPUHours

	// Operational Emissions are calculated by multiplying the energy, vCPUHours, PUE,
	// and region grid. The PUE is collected from the providers. The CO2e grid data
	// is the grid carbon intensity coefficient for the region at the specified time.
	p.metric.Emissions = v1.NewResourceEmission(
		p.metric.Energy*p.pue*p.grid,
		v1.GCO2eq,
	)

	logger.Debug("CPU calculation", "energy usage", p.metric.Energy, "emissions", p.metric.Emissions)
	return nil
}

// memory is calculated based on the TEADs pkgRAM calculations over
// various memory stress loads. Using the memory usage from the instance
// we can get the estimated Power consumption of the instance
func memory(ctx context.Context, p *parameters) error {
	logger := log.FromContext(ctx)
	var err error

	// if utilization is not found from the query, fallback to:
	// 1. Check if GB Ram usage is found (for GCP this is only
	//    collected for e2 instances) and divide by total GBs of
	//    memory allocated to VM to get the usage. If Ram used is
	//    not found,
	// 2. Set to 50% as a temporary placeholder until a more
	//    accurate solution is found.
	if p.metric.Usage == 0 {
		p.metric.Usage = 50
		if p.metric.Unit == v1.GB && p.metric.UnitAmount > 0 {
			p.metric.Usage = (p.metric.UnitAmount / p.factors.MemoryGB) * 100
		}
	}

	p.metric.Energy, err = cubicSplineInterpolation(p.factors.RAMWatt, p.metric.Usage)
	if err != nil {
		return err
	}

	p.metric.Emissions = v1.NewResourceEmission(
		p.metric.Energy*p.pue*p.grid,
		v1.GCO2eq,
	)

	logger.Debug("Memory calculation", "energy usage", p.metric.Energy, "emissions", p.metric.Emissions)
	return nil
}

// cubicSplineInterpolation is a piecewise cubic polynomials that takes the
// four measured wattage data points at 0%, 10%, 50%, and 100% utilization
// and interpolates a value for the usage (%) value and returns the energy
// in kilowatts.
func cubicSplineInterpolation(wattage []data.Wattage, value float64) (float64, error) {
	if len(wattage) == 0 {
		return 0, errors.New("error: cannot calculate CPU energy, no wattage found")
	}

	// split the wattage slice into a slice of
	// float percentages and a slice of wattages
	var x, y = []float64{}, []float64{}
	for _, w := range wattage {
		x = append(x, float64(w.Percentage))
		y = append(y, w.Wattage)
	}

	s := gospline.NewCubicSpline(x, y)
	// s.At returns the cubic spline value in Wattage
	// divide by 1000 to get kilowatts.
	return s.At(value) / 1000, nil
}

// EmbodiedEmissions are the released emissions of production and destruction of the
// hardware
func embodiedEmissions(interval time.Duration, hourlyEmbodied float64) float64 {
	// The embodied emissions need to be calculated for the measurement interval, so the
	// hourly emissions further divided to the interval minutes.
	return hourlyEmbodied / float64(60) * float64(interval.Minutes())
}
