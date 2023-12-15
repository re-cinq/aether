package v1

type CoefficientData map[string]float64       // map[region] = co2e
type EmbodiedData map[string]Embodied         // key = Machine type (n2-standard-
type MachineSpecsData map[string]MachineSpecs // key = architecture name (Haswell, Skylake, ..)

type EmissionFactors struct {
	Provider    string
	Coefficient CoefficientData // key is region
	Embodied    EmbodiedData    // key is machineType
	*ProviderDefaults
}

type Coefficient struct {
	Region string
	Co2e   float64
}

// TotalEmbodied assumes base manufacturing emissions of 1000 kgCO2e
// for a mono socket, low DRAM, no local storage AND
// the following additional manufacturing emissions:
// storage: 50 kgCO2 per HDD, 100 kgCO2 per SDD
// cpus: 100 kgCO2e per cpu
// gpus: 150 kgCO2e per gpu
// rack server lifespace: 4 years
// memory: (533 / 384) * DRAM_THRESHOLD
// Based on Dell PowerEdge R740 Life-Cycle Assessment
// = 533 kgCO₂eq for 12*32GB DIMMs Memory (384 GB).

// TODO: Need to add machine lifespan for embodied carbon
type Embodied struct {
	MachineType                   string  `yaml:"type"`
	AdditionalMemoryKiloWattCO2e  float64 `yaml:"additionalmemory"`
	AdditionalStorageKiloWattCO2e float64 `yaml:"additionalstorage"`
	AdditionalCPUsKiloWattCO2e    float64 `yaml:"additionalcpus"`
	AdditionalGPUsKiloWattCO2e    float64 `yaml:"additionalgpus"`
	TotalEmbodiedKiloWattCO2e     float64 `yaml:"total"`
	Architecture                  string
	MachineSpecs
}

type MachineSpecs struct {
	Architecture string
	MinWatts     float64
	MaxWatts     float64
	GBPerChip    float64 `yaml:"chip"`
}

type ProviderDefaults struct {
	Provider                 string  `yaml:"name"`
	MinWatts                 float64 `yaml:"minWatts"`
	MaxWatts                 float64 `yaml:"maxWatts"`
	HDDStorageWatts          float64 `yaml:"hddStorageWatts"`
	SSDStorageWatts          float64 `yaml:"ssdStorageWatts"`
	NetworkingKilloWattHours float64 `yaml:"networkingKilloWattHours"`
	MemoryKilloWattHours     float64 `yaml:"memoryKilloWattHours"`
	AveragePUE               float64 `yaml:"averagePUE"`
}
