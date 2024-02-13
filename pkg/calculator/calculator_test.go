package calculator

// import (
// 	"testing"
// 	"time"
//
// 	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
// 	"github.com/stretchr/testify/assert"
// )
//
// type testcase struct {
// 	name      string
// 	interval  time.Duration // this is nanoseconds
// 	calculate *calculate
// 	metric    *v1.Metric
// 	expRes    float64
// 	hasErr    bool
// 	expErr    string
// }
//
// // defaultCalc is contains basic/typical emissions
// // data numbers that are used as the default for tests
// func defaultCalc() *calculate {
// 	return &calculate{
// 		minWatts: 1.3423402398570,
// 		maxWatts: 4.00498247528,
// 		chip:     35.23458732,
// 		pue:      1.0123,
// 		gridCO2e: 0.00023,
// 	}
// }
//
// func defaultMetric() *v1.Metric {
// 	m := v1.NewMetric("basic")
// 	m.SetType(v1.CPU).SetUsage(25)
// 	m.SetUnitAmount(4).SetResourceUnit(v1.VCPU)
// 	return m
// }
//
// func TestCalculateCPUEmissions(t *testing.T) {
// 	for _, test := range []*testcase{
// 		func() *testcase {
// 			// Default test case
// 			return &testcase{
// 				name:      "basic default numbers",
// 				interval:  30 * time.Second,
// 				calculate: defaultCalc(),
// 				metric:    defaultMetric(),
// 				expRes:    0.0005270347987162735,
// 			}
// 		}(),
// 		func() *testcase {
// 			// All data set to zero values
// 			return &testcase{
// 				name:     "no values in calculator",
// 				interval: 30 * time.Second,
// 				calculate: &calculate{
// 					minWatts: 0,
// 					maxWatts: 0,
// 					chip:     0,
// 					pue:      0,
// 					gridCO2e: 0,
// 				},
// 				metric: defaultMetric(),
// 				expRes: 0,
// 			}
// 		}(),
// 		func() *testcase {
// 			// vCPUs not set
// 			return &testcase{
// 				name:      "no vCPUs set",
// 				interval:  30 * time.Second,
// 				calculate: defaultCalc(),
// 				metric:    defaultMetric().SetUnitAmount(0),
// 				hasErr:    true,
// 				expErr:    "error vCPUs set to 0, this should never be the case",
// 			}
// 		}(),
//
// 		func() *testcase {
// 			// Calculate the default values over
// 			// a 5 minute interval, instead of
// 			// 30 seconds
// 			return &testcase{
// 				name:      "5 minutes interval",
// 				interval:  5 * time.Minute,
// 				calculate: defaultCalc(),
// 				metric:    defaultMetric(),
// 				expRes:    0.005270347987162734,
// 			}
// 		}(),
//
// 		func() *testcase {
// 			// Calculate the default values over
// 			// one hour
// 			return &testcase{
// 				name:      "1 hour interval",
// 				interval:  1 * time.Hour,
// 				calculate: defaultCalc(),
// 				metric:    defaultMetric(),
// 				expRes:    0.06324417584595282,
// 			}
// 		}(),
//
// 		func() *testcase {
// 			// calculate with only a single vCPU
// 			return &testcase{
// 				name:      "single vCPU",
// 				interval:  30 * time.Second,
// 				calculate: defaultCalc(),
// 				metric:    defaultMetric().SetUnitAmount(1),
// 				expRes:    0.00013175869967906837,
// 			}
// 		}(),
//
// 		func() *testcase {
// 			// test with vCPU utilization at 50%
// 			return &testcase{
// 				name:      "50% utilization",
// 				interval:  30 * time.Second,
// 				calculate: defaultCalc(),
// 				metric:    defaultMetric().SetUsage(50),
// 				expRes:    0.0010436517395756915,
// 			}
// 		}(),
//
// 		func() *testcase {
// 			// test with vCPU utilization at 100%
// 			return &testcase{
// 				name:      "100% utilization",
// 				interval:  30 * time.Second,
// 				calculate: defaultCalc(),
// 				metric:    defaultMetric().SetUsage(100),
// 				expRes:    0.002076885621294527,
// 			}
// 		}(),
//
// 		func() *testcase {
// 			c := defaultCalc()
// 			c.pue = 1.0
// 			// test if PUE is exactly 1
// 			return &testcase{
// 				name:      "PUE is exactly 1.0",
// 				interval:  30 * time.Second,
// 				calculate: c,
// 				metric:    defaultMetric(),
// 				expRes:    0.0005206310369616452,
// 			}
// 		}(),
//
// 		func() *testcase {
// 			c := defaultCalc()
// 			c.gridCO2e = 402
// 			// test an extremely high grid CO2e
// 			// This value was collected from azures
// 			// Germany West Central region
// 			return &testcase{
// 				name:      "High grid CO2e",
// 				interval:  30 * time.Second,
// 				calculate: c,
// 				metric:    defaultMetric(),
// 				expRes:    921.1651699301823,
// 			}
// 		}(),
//
// 		func() *testcase {
// 			// create a relatively large server with higher
// 			// than typical min and max watts, 32 vCPUs, and
// 			// a utilization of 90%
// 			return &testcase{
// 				name:     "large server and large workload",
// 				interval: 30 * time.Second,
// 				calculate: &calculate{
// 					minWatts: 3.0369270833333335,
// 					maxWatts: 8.575357663690477,
// 					chip:     129.77777777777777,
// 					pue:      1.1,
// 					gridCO2e: 0.00079,
// 				},
// 				metric: defaultMetric().SetUnitAmount(32).SetUsage(90),
// 				expRes: 0.11621326542003971,
// 			}
// 		}(),
// 	} {
// 		t.Run(test.name, func(t *testing.T) {
// 			res, err := test.calculate.cpu(test.metric, test.interval)
// 			assert.Equalf(t, test.expRes, res, "Result should be: %v, got: %v", test.expRes, res)
// 			if test.hasErr {
// 				assert.EqualErrorf(t, err, test.expErr, "Error should be: %v, got: %v", test.expErr, err)
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 		})
// 	}
// }
//
// func TestCalculateEmbodiedEmissions(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		embodied float64
// 		expRes   float64
// 		interval time.Duration
// 	}{
// 		{
// 			name:     "1706.48 total 30 seconds",
// 			embodied: 1706.48,
// 			expRes:   0.000270560629122273,
// 			interval: 30 * time.Second,
// 		},
// 		{
// 			name:     "1706.48 total 5 minutes",
// 			embodied: 1706.48,
// 			expRes:   0.0027056062912227297,
// 			interval: 5 * time.Minute,
// 		},
// 		{
// 			name:     "no embodied emissions",
// 			embodied: 0,
// 			expRes:   0,
// 			interval: 30 * time.Second,
// 		},
// 		{
// 			name:     "large emissions, 1 minute",
// 			embodied: 6268.55,
// 			expRes:   0.0019877441653982754,
// 			interval: 1 * time.Minute,
// 		},
// 	}
//
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			c := calculate{totalEmbodied: test.embodied}
// 			res := c.embodiedEmissions(test.interval)
// 			assert.Equalf(t, test.expRes, res, "Result should be: %v, got: %v", test.expRes, res)
// 		})
// 	}
// }
