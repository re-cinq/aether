package v1

import (
	"fmt"
	"testing"

	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	"github.com/stretchr/testify/assert"
)

const testDataPath = "../../../testdata"

func TestGetProviderDefaults(t *testing.T) {
	tests := []struct {
		name     string
		provider v1.Provider
		hasError bool
		expRes   *ProviderDefaults
		expErr   string
	}{
		{
			name:     "pass: read and set provider defaults",
			provider: "fake",
			hasError: false,
			expRes: &ProviderDefaults{
				Provider:                 "fake",
				MinWatts:                 0.71,
				MaxWatts:                 3.5,
				HDDStorageWatts:          0.65,
				SSDStorageWatts:          1.22,
				NetworkingKilloWattHours: 0.001,
				MemoryKilloWattHours:     0.000392,
				AveragePUE:               1.125,
			},
			expErr: "",
		},
		{
			name:     "fail: fails to read data",
			provider: "bad",
			hasError: true,
			expRes:   nil,
			expErr:   fmt.Sprintf("open %s/bad-default.yaml: no such file or directory", testDataPath),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ef := &EmissionFactors{Provider: test.provider}
			err := ef.getProviderDefaults(testDataPath)
			assert.Equalf(t, ef.ProviderDefaults, test.expRes, "Result should be: %v, got: %v", test.expRes, ef.ProviderDefaults)
			if test.hasError {
				assert.EqualErrorf(t, err, test.expErr, "Error should be: %v, got: %v", test.expErr, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetCoefficientData(t *testing.T) {
	tests := []struct {
		name     string
		provider v1.Provider
		hasError bool
		expRes   CoefficientData
		expErr   string
	}{
		{
			name:     "pass: create mapping of coefficient data",
			provider: "fake",
			hasError: false,
			expRes: CoefficientData{
				"us-central1":     0.000479,
				"us-east1":        0.0005,
				"ap-northeast-1":  0.000506,
				"ca-central-1":    0.00013,
				"France Central":  6.7e-05,
				"Finland Central": 77,
			},
			expErr: "",
		},
		{
			name:     "fail: malformed YAML file",
			provider: "bad",
			hasError: true,
			expRes:   CoefficientData{},
			expErr:   "yaml: line 1: mapping values are not allowed in this context",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ef := &EmissionFactors{Provider: test.provider}
			err := ef.getCoefficientData(testDataPath)
			assert.Equalf(t, ef.Coefficient, test.expRes, "Result should be: %v, got: %v", test.expRes, ef.Coefficient)
			if test.hasError {
				assert.EqualErrorf(t, err, test.expErr, "Error should be: %v, got: %v", test.expErr, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetMachineSpec(t *testing.T) {
	tests := []struct {
		name     string
		provider v1.Provider
		hasError bool
		expRes   MachineSpecsData
		expErr   string
	}{
		{
			name:     "pass: read and create machineSpec dataset",
			provider: "fake",
			hasError: false,
			expRes: MachineSpecsData{
				"Broadwell": {
					Architecture: "Broadwell",
					MinWatts:     0.7128342245989304,
					MaxWatts:     3.3857473048128344,
					GBPerChip:    69.6470588235294,
				},
				"Haswell": {
					Architecture: "Haswell",
					MinWatts:     1.9005681818181814,
					MaxWatts:     5.9688982156043195,
					GBPerChip:    27.310344827586206,
				},
				"Skylake": {
					Architecture: "Skylake",
					MinWatts:     0.6446044454253452,
					MaxWatts:     3.8984738056304855,
					GBPerChip:    80.43037974683544,
				},
				"EPYC 2nd Gen": {
					Architecture: "EPYC 2nd Gen",
					MinWatts:     0.4742621527777778,
					MaxWatts:     1.5751872939814815,
					GBPerChip:    129.77777777777777,
				},
			},
			expErr: "",
		},
		{
			name:     "fail: fails to read data",
			provider: "bad",
			hasError: true,
			expRes:   nil,
			expErr:   fmt.Sprintf("open %s/bad-use.yaml: no such file or directory", testDataPath),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := getMachineSpecs(test.provider, testDataPath)
			assert.Equalf(t, res, test.expRes, "Result should be: %v, got: %v", test.expRes, res)
			if test.hasError {
				assert.EqualErrorf(t, err, test.expErr, "Error should be: %v, got: %v", test.expErr, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetEmbodiedData(t *testing.T) {
	tests := []struct {
		name     string
		ef       EmissionFactors
		hasError bool
		expRes   EmbodiedData
		expErr   string
	}{
		{
			name:     "pass: read and set embodied with architecture",
			ef:       EmissionFactors{Provider: "fake"},
			hasError: false,
			expRes: EmbodiedData{
				"e2-standard-2": {
					MachineType:                   "e2-standard-2",
					AdditionalMemoryKiloWattCO2e:  155.46,
					AdditionalStorageKiloWattCO2e: 0,
					AdditionalCPUsKiloWattCO2e:    100,
					AdditionalGPUsKiloWattCO2e:    0,
					TotalEmbodiedKiloWattCO2e:     1255.46,
					Architecture:                  "Skylake",
					MachineSpecs: MachineSpecs{
						Architecture: "Skylake",
						MinWatts:     0.6446044454253452,
						MaxWatts:     3.8984738056304855,
						GBPerChip:    80.43037974683544,
					},
				},
				"n1-standard-2": {
					MachineType:                   "n1-standard-2",
					AdditionalMemoryKiloWattCO2e:  477.48,
					AdditionalStorageKiloWattCO2e: 100,
					AdditionalCPUsKiloWattCO2e:    100,
					AdditionalGPUsKiloWattCO2e:    0,
					TotalEmbodiedKiloWattCO2e:     1677.48,
					Architecture:                  "Broadwell",
					MachineSpecs: MachineSpecs{
						Architecture: "Broadwell",
						MinWatts:     0.7128342245989304,
						MaxWatts:     3.3857473048128344,
						GBPerChip:    69.6470588235294,
					},
				},
			},
			expErr: "",
		},
		{
			name:     "fail: fails to read data for machine specs",
			ef:       EmissionFactors{Provider: "bad"},
			hasError: true,
			expRes:   EmbodiedData{},
			expErr:   fmt.Sprintf("open %s/bad-use.yaml: no such file or directory", testDataPath),
		},
		{
			name:     "fail: architecture unknown and no provider defaults",
			ef:       EmissionFactors{Provider: "fake2"},
			hasError: true,
			expRes:   EmbodiedData{},
			expErr:   "error: machine specifications for architecture () does not exist nor do provider defaults",
		},
		{
			name: "pass: architecture unknown use provider Defaults",
			ef: EmissionFactors{
				Provider: "fake2",
				ProviderDefaults: &ProviderDefaults{
					MinWatts:   0.73,
					MaxWatts:   4.06,
					AveragePUE: 1.2,
				},
			},
			hasError: false,
			expRes: EmbodiedData{
				"c1.medium": {
					MachineType:                   "c1.medium",
					AdditionalMemoryKiloWattCO2e:  36.09,
					AdditionalStorageKiloWattCO2e: 400,
					AdditionalCPUsKiloWattCO2e:    100,
					AdditionalGPUsKiloWattCO2e:    0,
					TotalEmbodiedKiloWattCO2e:     1536.09,
					Architecture:                  "",
					MachineSpecs: MachineSpecs{
						Architecture: "",
						MinWatts:     0.73,
						MaxWatts:     4.06,
						GBPerChip:    0,
					},
				},
				"B1MS": {
					MachineType:                   "B1MS",
					AdditionalMemoryKiloWattCO2e:  88.83,
					AdditionalStorageKiloWattCO2e: 50,
					AdditionalCPUsKiloWattCO2e:    100,
					AdditionalGPUsKiloWattCO2e:    0,
					TotalEmbodiedKiloWattCO2e:     1238.83,
					Architecture:                  "Unknown",
					MachineSpecs: MachineSpecs{
						Architecture: "",
						MinWatts:     0.73,
						MaxWatts:     4.06,
						GBPerChip:    0,
					},
				},
			},
			expErr: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.ef.getEmbodiedData(testDataPath)
			assert.Equalf(t, test.ef.Embodied, test.expRes, "Result should be: %v, got: %v", test.expRes, test.ef.Embodied)
			if test.hasError {
				assert.EqualErrorf(t, err, test.expErr, "Error should be: %v, got: %v", test.expErr, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetEmissionFactors(t *testing.T) {
	tests := []struct {
		name     string
		provider v1.Provider
		hasError bool
		expRes   *EmissionFactors
		expErr   string
	}{
		{
			name:     "pass: read and create EmissionFactors",
			provider: "fake",
			hasError: false,
			expRes: &EmissionFactors{
				Provider: "fake",
				ProviderDefaults: &ProviderDefaults{
					Provider:                 "fake",
					MinWatts:                 0.71,
					MaxWatts:                 3.5,
					HDDStorageWatts:          0.65,
					SSDStorageWatts:          1.22,
					NetworkingKilloWattHours: 0.001,
					MemoryKilloWattHours:     0.000392,
					AveragePUE:               1.125,
				},
				Coefficient: CoefficientData{
					"us-central1":     0.000479,
					"us-east1":        0.0005,
					"ap-northeast-1":  0.000506,
					"ca-central-1":    0.00013,
					"France Central":  6.7e-05,
					"Finland Central": 77,
				},
				Embodied: EmbodiedData{
					"e2-standard-2": {
						MachineType:                   "e2-standard-2",
						AdditionalMemoryKiloWattCO2e:  155.46,
						AdditionalStorageKiloWattCO2e: 0,
						AdditionalCPUsKiloWattCO2e:    100,
						AdditionalGPUsKiloWattCO2e:    0,
						TotalEmbodiedKiloWattCO2e:     1255.46,
						Architecture:                  "Skylake",
						MachineSpecs: MachineSpecs{
							Architecture: "Skylake",
							MinWatts:     0.6446044454253452,
							MaxWatts:     3.8984738056304855,
							GBPerChip:    80.43037974683544,
						},
					},
					"n1-standard-2": {
						MachineType:                   "n1-standard-2",
						AdditionalMemoryKiloWattCO2e:  477.48,
						AdditionalStorageKiloWattCO2e: 100,
						AdditionalCPUsKiloWattCO2e:    100,
						AdditionalGPUsKiloWattCO2e:    0,
						TotalEmbodiedKiloWattCO2e:     1677.48,
						Architecture:                  "Broadwell",
						MachineSpecs: MachineSpecs{
							Architecture: "Broadwell",
							MinWatts:     0.7128342245989304,
							MaxWatts:     3.3857473048128344,
							GBPerChip:    69.6470588235294,
						},
					},
				},
			},
			expErr: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := GetEmissionFactors(test.provider, testDataPath)
			assert.Equalf(t, res, test.expRes, "Result should be: %v, got: %v", test.expRes, res)
			if test.hasError {
				assert.EqualErrorf(t, err, test.expErr, "Error should be: %v, got: %v", test.expErr, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
