package calculator

import (
	"context"
	"testing"
	"time"

	v1 "github.com/re-cinq/aether/pkg/types/v1"
	data "github.com/re-cinq/emissions-data/pkg/types/v2"
	"github.com/stretchr/testify/assert"
)

func params() *parameters {
	m := &v1.Metric{
		Name:  "basic",
		Usage: 27,
	}
	return &parameters{
		grid:   7,
		pue:    1.2,
		metric: m,
		// using t3.micro AWS instance as default
		factors: &data.Instance{
			PkgWatt: []data.Wattage{
				{
					Percentage: 0,
					Wattage:    1.21,
				},
				{
					Percentage: 10,
					Wattage:    3.05,
				},
				{
					Percentage: 50,
					Wattage:    7.16,
				},
				{
					Percentage: 100,
					Wattage:    9.96,
				},
			},
			RAMWatt: []data.Wattage{
				{
					Percentage: 0,
					Wattage:    0.15,
				},
				{
					Percentage: 10,
					Wattage:    0.24,
				},
				{
					Percentage: 50,
					Wattage:    0.62,
				},
				{
					Percentage: 100,
					Wattage:    1.00,
				},
			},
			VCPU: 2.0,
		},
		embodiedFactor: 1000,
	}
}

func TestCalculateCPU(t *testing.T) {
	type testcase struct {
		name      string
		interval  time.Duration // this is nanoseconds
		params    *parameters
		energy    float64
		emissions v1.ResourceEmissions
		hasErr    bool
		expErr    string
	}

	for _, test := range []*testcase{
		func() *testcase {
			// Default test case
			return &testcase{
				name:     "default t3.micro at 27% usage over 5m",
				interval: 5 * time.Minute,
				params:   params(),
				emissions: v1.ResourceEmissions{
					Value: 0.007453764094512195,
					Unit:  v1.GCO2eq,
				},
				energy: 0.0008873528683943089,
			}
		}(),
		func() *testcase {
			// vCPUs not set in params, but set in metric
			p := params()
			p.factors.VCPU = 0
			p.metric.UnitAmount = 2
			p.metric.Unit = v1.VCPU
			return &testcase{
				name:     "vCPU set in metric, not params",
				interval: 5 * time.Minute,
				params:   p,
				emissions: v1.ResourceEmissions{
					Value: 0.007453764094512195,
					Unit:  v1.GCO2eq,
				},
				energy: 0.0008873528683943089,
			}
		}(),
		func() *testcase {
			// vCPUs not set
			p := params()
			p.factors.VCPU = 0
			return &testcase{
				name:     "vCPU not set",
				interval: 5 * time.Minute,
				params:   p,
				hasErr:   true,
				expErr:   "error vCPU set to 0",
				emissions: v1.ResourceEmissions{
					Value: 0.0,
					Unit:  "",
				},
				energy: 0.0,
			}
		}(),

		func() *testcase {
			// Calculate the default values over
			// a 30 second interval, instead of
			// 5 minutes
			return &testcase{
				name:     "default 30 second interval",
				interval: 30 * time.Second,
				params:   params(),
				emissions: v1.ResourceEmissions{
					Value: 0.0007453764094512195,
					Unit:  v1.GCO2eq,
				},
				energy: 8.87352868394309e-05,
			}
		}(),

		func() *testcase {
			// Calculate the default values over
			// one hour
			return &testcase{
				name:     "1 hour interval",
				interval: 1 * time.Hour,
				params:   params(),
				emissions: v1.ResourceEmissions{
					Value: 0.08944516913414635,
					Unit:  v1.GCO2eq,
				},
				energy: 0.010648234420731708,
			}
		}(),

		func() *testcase {
			// calculate with 4 vCPUs
			p := params()
			p.factors.VCPU = 4
			return &testcase{
				name:     "4 vCPU",
				interval: 5 * time.Minute,
				params:   p,
				emissions: v1.ResourceEmissions{
					Value: 0.01490752818902439,
					Unit:  v1.GCO2eq,
				},
				energy: 0.0017747057367886179,
			}
		}(),

		func() *testcase {
			p := params()
			p.pue = 1.0
			// test if PUE is exactly 1
			return &testcase{
				name:     "PUE is exactly 1.0",
				interval: 5 * time.Minute,
				params:   p,
				emissions: v1.ResourceEmissions{
					Value: 0.006211470078760163,
					Unit:  v1.GCO2eq,
				},
				energy: 0.0008873528683943089,
			}
		}(),

		func() *testcase {
			p := params()
			p.grid = 402
			// test an extremely high grid CO2e
			// This value was collected from azures
			// Germany West Central region
			return &testcase{
				name:     "High grid CO2e",
				interval: 5 * time.Minute,
				params:   p,
				emissions: v1.ResourceEmissions{
					Value: 0.42805902371341464,
					Unit:  v1.GCO2eq,
				},
				energy: 0.0008873528683943089,
			}
		}(),

		func() *testcase {
			// create a relatively large server with higher
			// than typical min and max watts, 32 vCPUs, and
			// a utilization of 90%
			p := params()
			p.metric.UnitAmount = 32
			p.metric.Usage = 90
			return &testcase{
				name:     "large server and large workload",
				interval: 5 * time.Minute,
				params:   p,
				emissions: v1.ResourceEmissions{
					Value: 0.013218390243902438,
					Unit:  v1.GCO2eq,
				},
				energy: 0.0015736178861788619,
			}
		}(),
	} {
		t.Run(test.name, func(t *testing.T) {
			err := cpu(context.TODO(), test.interval, test.params)
			actualEnergy := test.params.metric.Energy
			actualEmissions := test.params.metric.Emissions

			assert.Equalf(t, test.energy, actualEnergy, "Result should be: %v, got: %v", test.energy, actualEnergy)
			assert.Equalf(t, test.emissions, actualEmissions, "Result should be: %v, got: %v", test.emissions, actualEmissions)
			if test.hasErr {
				assert.Errorf(t, err, test.expErr)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestCubicSplineInterpolation(t *testing.T) {
	type testcase struct {
		name      string
		powerCPU  []data.Wattage
		usage     float64
		emissions float64
		hasErr    bool
		expErr    string
	}
	testcases := []testcase{
		{
			name: "t3.micro instance at 27%",
			powerCPU: []data.Wattage{
				{
					Percentage: 0,
					Wattage:    1.21,
				},
				{
					Percentage: 10,
					Wattage:    3.05,
				},
				{
					Percentage: 50,
					Wattage:    7.16,
				},
				{
					Percentage: 100,
					Wattage:    9.96,
				},
			},
			usage:     27.00,
			emissions: 0.005324117210365854,
		},
		{
			name:      "empty wattage",
			powerCPU:  []data.Wattage{},
			usage:     27.01,
			emissions: 0,
			hasErr:    true,
			expErr:    "error: cannot calculate CPU energy, no wattage found",
		},
		{
			name: "at exactly 10% utilization",
			powerCPU: []data.Wattage{
				{
					Percentage: 0,
					Wattage:    1.21,
				},
				{
					Percentage: 10,
					Wattage:    3.05,
				},
				{
					Percentage: 50,
					Wattage:    7.16,
				},
				{
					Percentage: 100,
					Wattage:    9.96,
				},
			},
			usage:     10,
			emissions: 0.0030499999999999993,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			res, err := cubicSplineInterpolation(test.powerCPU, test.usage)
			if test.hasErr {
				assert.Errorf(t, err, test.expErr)
			} else {
				assert.Nil(t, err)
			}
			assert.Equalf(t, test.emissions, res, "Result should be: %v, got: %v", test.emissions, res)
		})
	}
}

func TestCalculateMemory(t *testing.T) {
	type testcase struct {
		name      string
		params    *parameters
		emissions float64
		energy    float64
		hasErr    bool
		expErr    string
	}
	for _, test := range []*testcase{
		func() *testcase {
			// pass: default test case
			return &testcase{
				name:      "default t3.micro at 27%",
				params:    params(),
				energy:    0.00040240120731707316,
				emissions: 0.0033801701414634144,
			}
		}(),
		func() *testcase {
			// fail: powerRAM wattage not set
			p := params()
			p.factors.RAMWatt = []data.Wattage{}
			return &testcase{
				name:      "fail: wattage RAM data not set",
				params:    p,
				energy:    0,
				emissions: 0,
				hasErr:    true,
				expErr:    "RAM wattage data not found for memory calculation",
			}
		}(),
	} {
		t.Run(test.name, func(t *testing.T) {
			err := memory(context.TODO(), test.params)
			actualEmissions := test.params.metric.Emissions.Value
			actualEnergy := test.params.metric.Energy

			assert.Equalf(t, test.energy, actualEnergy, "Result should be: %v, got: %v", actualEnergy, test.energy)
			assert.Equalf(t, test.emissions, actualEmissions, "Result should be: %v, got: %v", actualEmissions, test.emissions)
			if test.hasErr {
				assert.Errorf(t, err, test.expErr)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
