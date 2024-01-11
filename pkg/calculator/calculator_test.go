package calculator

import (
	"testing"
	"time"

	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	"github.com/stretchr/testify/assert"
)

type testcase struct {
	name     string
	interval time.Duration // this is nanoseconds
	calc     *calculate
	expRes   float64
}

func basicTestcase() *testcase {
	return &testcase{
		name:     "basic",
		interval: 30 * time.Second,
		calc: &calculate{
			cores:    4,
			usageCPU: v1.Percentage(25),
			minWatts: 1.3423402398570,
			maxWatts: 4.00498247528,
			chip:     35.23458732,
			pue:      1.0123,
			gridCO2e: 0.00023,
		},
		expRes: 0.0005270347987162735,
	}
}

func noValues() *testcase {
	return &testcase{
		name:     "no values 30 sec",
		interval: 30 * time.Second,
		calc: &calculate{
			cores:    0,
			usageCPU: v1.Percentage(0),
			minWatts: 0,
			maxWatts: 0,
			chip:     0,
			pue:      0,
			gridCO2e: 0,
		},
		expRes: 0,
	}
}

func fiveMinuteInterval() *testcase {
	return &testcase{
		name:     "5 minutes interval",
		interval: 5 * time.Minute,
		calc:     basicTestcase().calc,
		expRes:   0.005270347987162734,
	}
}

func oneHourInterval() *testcase {
	return &testcase{
		name:     "5 minutes interval",
		interval: 1 * time.Hour,
		calc:     basicTestcase().calc,
		expRes:   0.06324417584595282,
	}
}

func singleCore() *testcase {
	c := basicTestcase().calc
	c.cores = 1
	return &testcase{
		name:     "single core",
		interval: 30 * time.Second,
		calc:     c,
		expRes:   0.00013175869967906837,
	}
}

func fiftyPercentCPUUtilization() *testcase {
	c := basicTestcase().calc
	c.usageCPU = v1.Percentage(50)
	return &testcase{
		name:     "50% utilization",
		interval: 30 * time.Second,
		calc:     c,
		expRes:   0.0010436517395756915,
	}
}

func OneHundredPercentCPUUtilization() *testcase {
	c := basicTestcase().calc
	c.usageCPU = v1.Percentage(100)
	return &testcase{
		name:     "100% utilization",
		interval: 30 * time.Second,
		calc:     c,
		expRes:   0.002076885621294527,
	}
}

func highPUE() *testcase {
	c := basicTestcase().calc
	c.pue = 1.777
	return &testcase{
		name:     "high PUE",
		interval: 30 * time.Second,
		calc:     c,
		expRes:   0.0009251613526808435,
	}
}

func highGridCO2e() *testcase {
	c := basicTestcase().calc
	c.gridCO2e = 5.2e-05
	return &testcase{
		name:     "High grid CO2e",
		interval: 30 * time.Second,
		calc:     c,
		expRes:   0.00011915569362280964,
	}
}

func largeServerLargeWorkload() *testcase {
	return &testcase{
		name:     "large server and large workload",
		interval: 30 * time.Second,
		calc: &calculate{
			cores:    32,
			usageCPU: v1.Percentage(90),
			minWatts: 3.0369270833333335,
			maxWatts: 8.575357663690477,
			chip:     129.77777777777777,
			pue:      1.1,
			gridCO2e: 0.00079,
		},
		expRes: 0.11621326542003971,
	}
}
func TestCalculateEmissions(t *testing.T) {
	testcases := []*testcase{
		basicTestcase(),
		noValues(),
		fiveMinuteInterval(),
		oneHourInterval(),
		singleCore(),
		fiftyPercentCPUUtilization(),
		OneHundredPercentCPUUtilization(),
		highPUE(),
		highGridCO2e(),
		largeServerLargeWorkload(),
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			res := test.calc.operationalCPUEmissions(test.interval)
			assert.Equalf(t, test.expRes, res, "Result should be: %v, got: %v", test.expRes, res)
		})
	}
}

func TestCalculateEmbodiedEmissions(t *testing.T) {
	tests := []struct {
		name     string
		embodied float64
		expRes   float64
		interval time.Duration
	}{
		{
			name:     "1706.48 total 30 seconds",
			embodied: 1706.48,
			expRes:   0.0004058409436834094,
			interval: 30 * time.Second,
		},
		{
			name:     "1706.48 total 5 minutes",
			embodied: 1706.48,
			expRes:   0.004058409436834094,
			interval: 5 * time.Minute,
		},
		{
			name:     "no embodied emissions",
			embodied: 0,
			expRes:   0,
			interval: 30 * time.Second,
		},
		{
			name:     "large emissions, 1 minute",
			embodied: 6268.55,
			expRes:   0.0029816162480974123,
			interval: 1 * time.Minute,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := calculate{totalEmbodied: test.embodied}
			res := c.embodiedEmissions(test.interval)
			assert.Equalf(t, test.expRes, res, "Result should be: %v, got: %v", test.expRes, res)
		})
	}
}
