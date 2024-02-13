package calculator

import (
	"errors"
	"fmt"
	"time"

	"github.com/cnkei/gospline"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	data "github.com/re-cinq/emissions-data/pkg/types/v2"
	"k8s.io/klog"
)

type parameters struct {
	gridCO2e float64
	pue      float64
	metric   v1.Metric
	dataset  data.Instance
}

// AWS, GCP and Azure have increased their server lifespan to 6 years (2024)
// https://sustainability.aboutamazon.com/products-services/the-cloud?energyType=true
// https://www.theregister.com/2024/01/31/alphabet_q4_2023/#:~:text=Alphabet%20first%20decided%20to%20extend,for%20six%20years%20before%20replacement.
// https://www.theregister.com/2022/08/02/microsoft_server_life_extension/
const serverLifespan = 6

// operationalEmissions determines the correct function to run to calculate the
// operational emissions for the metric type
func operationalEmissions(interval time.Duration, p *parameters) (float64, error) {
	switch p.metric.Name() {
	case v1.CPU.String():
		return cpu(interval, p)
	case v1.Memory.String():
		return memory(interval, p)
	case v1.Storage.String():
		return 0, errors.New("error storage is not yet being calculated")
	case v1.Network.String():
		return 0, errors.New("error networking is not yet being calculated")
	default:
		return 0, fmt.Errorf("error metric not supported: %+v", p.metric.Name())
	}
}

// cpu calculates the CO2e operational emissions for the CPU utilization of
// a Cloud VM instance over an interval of time.
// More information of the calculation can be found in the docs.
//
// The initial calculation uses the wattage conversion factor based on the turbostat and
// turbostress to stress test the CPU on baremetal servers as inspired by Teads.
// If those datasets do not exist, we fall back to calculate based on min and max
// wattage from SPECpower data as inspired by CCF and Etsy.
func cpu(interval time.Duration, p *parameters) (float64, error) {
	vCPU := float64(p.dataset.VCPU)
	// vCPU are virtual CPUs that are mapped to physical cores (a core is a physical
	// component to the CPU the VM is running on). If vCPU from the dataset (p.vCPU)
	// is not found, get the number of vCPUs from the metric collected from the query
	if vCPU == 0 {
		if p.metric.UnitAmount() == 0 {
			return 0, errors.New("error vCPU set to 0")
		}
		vCPU = p.metric.UnitAmount()
	}

	// Calculate CPU energy consumption as a rate by the interval time for 1 hour.
	// For example, if the machine has 4 vCPUs and the interval of time is
	// 5 minutes: The hourly time is 5/60 (0.083333333) * 4 vCPU  = 0.33333334
	vCPUHours := (interval.Minutes() / float64(60)) * vCPU

	// if there pkgWatt dataset values, then use interpolation
	// to calculate the wattage based on the utilization, otherwise, calculate
	// based on SPECpower min and max data
	fmt.Println("prespline: ", p.dataset.PkgWatt, p.metric.Usage())
	watts := cubicSplineInterpolation(p.dataset.PkgWatt, p.metric.Usage())
	klog.Infof("CPU interpolation result: %+v from wattage: %+v and usage %+v\n", watts, p.dataset.PkgWatt, p.metric.Usage())
	if len(p.dataset.PkgWatt) == 2 {
		// Average Watts is the average energy consumption of the service. It is based on
		// CPU utilization and Minimum and Maximum wattage of the server. If the machine
		// architecture is unknown the Min and Max wattage is the average of all machines
		// for that provider, and is supplied in the provider defaults. This is being
		// handled in the types/factors package (the point of reading in coefficient data).
		minWatts := p.dataset.PkgWatt[0].Wattage
		maxWatts := p.dataset.PkgWatt[1].Wattage
		watts = minWatts + p.metric.Usage()*(maxWatts-minWatts)
	}
	// Operational Emissions are calculated by multiplying the avgWatts, vCPUHours, PUE,
	// and region grid CO2e. The PUE is collected from the providers. The CO2e grid data
	// is the electrical grid emissions for the region at the specified time.
	return watts * vCPUHours * p.pue * p.gridCO2e, nil
}

func memory(interval time.Duration, p *parameters) (float64, error) {
	fmt.Printf("memory: %+v\n", p.dataset)
	if len(p.dataset.RAMWatt) == 0 {
		return 0, nil
	}

	coefficient := 0.000392 // kw/GB hrs

	GBHours := p.metric.UnitAmount() * (interval.Minutes() / float64(60)) * coefficient

	watts := cubicSplineInterpolation(p.dataset.RAMWatt, p.metric.Usage())
	// amount GB * interval?
	kwHours := watts * GBHours * p.pue * p.gridCO2e
	fmt.Printf("XX: %+v %+v\n", GBHours, kwHours)
	fmt.Printf("Memory metric: %+v\n", p.metric)
	return kwHours, nil
}

// TODO add comment
func cubicSplineInterpolation(wattage []data.Wattage, value float64) float64 {
	fmt.Println("wattage: ", wattage)
	var x, y = []float64{}, []float64{}
	for _, w := range wattage {
		x = append(x, float64(w.Percentage))
		y = append(y, w.Wattage)
	}

	s := gospline.NewCubicSpline(x, y)
	return s.At(value)
}

// EmbodiedEmissions are the released emissions of production and destruction of the
// hardware
func embodiedEmissions(interval time.Duration, totalEmbodied float64) float64 {
	// Total Embodied is the total emissions for a server to be produced, including
	// additional emmissions for added DRAM, CPUs, GPUS, and storage. This is divided
	// by the expected lifespan of the server to get the annual emissions.
	annualEmbodied := totalEmbodied / serverLifespan

	// The embodied emissions need to be calculated for the measurement interval, so the
	// annual emissions further divided to the interval minutes.
	//nolint:unconvert //conversion to minutes does affect calculation
	return annualEmbodied / float64(365) / float64(24) / float64(60) * float64(interval.Minutes())
}
