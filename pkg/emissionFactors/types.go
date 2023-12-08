package emissionfactors

type emissionFactors struct {
	provider         string
	providerDefaults defaults      // not sure about this
	coefficient      []coefficient // region dependent
	embodied         embodied
}

type coefficient struct {
	region string
	co2e   float64
}

type embodied struct {
	machineType       string `yaml:"type"`
	additionalmemory  float64
	additionalstorage float64
	additionalcpus    float64
	additionalgpus    float64
	total             float64
	architecture      architecture
}

type architecture struct {
	name     string
	minWatts float64
	maxWatts float64
	chip     float64
}

type defaults struct {
	provider                 string `yaml:"name"`
	minWatts                 float64
	maxWatts                 float64
	hddStorageWatts          float64
	ssdStorageWatts          float64
	networkingKilloWattHours float64 // per Gigabyte
	memoryKilloWattHours     float64 // per Gigabyte
}
