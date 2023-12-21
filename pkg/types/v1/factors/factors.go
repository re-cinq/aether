package v1

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
)

// Emission data is currently stored as files in our emissions-data repo.
// Each file is named "{provider}-{emissionFactor}" where emissionFactor
// may be default, embodied, grid, and use.

// GetEmissionFactors reads in emission data from files and stores them
// into the emissionFactors struct for calulating

// TODO put this is in config
// dataPath := "../../../emissions-data/data/"
func GetEmissionFactors(provider v1.Provider, dataPath string) (*EmissionFactors, error) {
	var err error
	ef := &EmissionFactors{
		Provider: provider,
	}

	err = ef.getProviderDefaults(dataPath)
	if err != nil {
		return nil, err
	}

	err = ef.getCoefficientData(dataPath)
	if err != nil {
		return nil, err
	}

	err = ef.getEmbodiedData(dataPath)
	if err != nil {
		return nil, err
	}

	return ef, nil
}

func (ef *EmissionFactors) getProviderDefaults(dataPath string) error {
	data := &ProviderDefaults{}

	fp := filepath.Join(dataPath, fmt.Sprintf("%s-default.yaml", ef.Provider))
	err := readYamlData(fp, data)
	if err != nil {
		return err
	}

	ef.ProviderDefaults = data
	return nil
}

// getCoefficeintData reads the {provider}-grid.yaml file into a slice
// of Coefficient structs, and then converts the data into a map of
// region: co2e to be returned
func (ef *EmissionFactors) getCoefficientData(dataPath string) error {
	data := []Coefficient{}
	ef.Coefficient = make(CoefficientData)

	fp := filepath.Join(dataPath, fmt.Sprintf("%s-grid.yaml", ef.Provider))
	if err := readYamlData(fp, &data); err != nil {
		return err
	}

	for _, c := range data {
		ef.Coefficient[c.Region] = c.Co2e
	}

	return nil
}

// getMachineSpecs creates a map of machine specs based on
// machine architecture
func getMachineSpecs(provider v1.Provider, dataPath string) (MachineSpecsData, error) {
	data := []MachineSpecs{}
	machineSpecsData := make(MachineSpecsData)

	fp := filepath.Join(dataPath, fmt.Sprintf("%s-use.yaml", provider))
	err := readYamlData(fp, &data)
	if err != nil {
		return nil, err
	}

	for _, d := range data {
		machineSpecsData[d.Architecture] = d
	}

	return machineSpecsData, nil
}

// getEmbodiedData maps an embedded struct of machineSpecs into
// the Embodied struct by the machine type
func (ef *EmissionFactors) getEmbodiedData(dataPath string) error {
	data := []Embodied{}
	ef.Embodied = make(EmbodiedData)

	machineSpecsData, err := getMachineSpecs(ef.Provider, dataPath)
	if err != nil {
		return err
	}

	fp := filepath.Join(dataPath, fmt.Sprintf("%s-embodied.yaml", ef.Provider))
	err = readYamlData(fp, &data)
	if err != nil {
		return err
	}

	for _, d := range data {
		val, ok := machineSpecsData[d.Architecture]
		// use provider defaults if architecture cannot be found
		if !ok {
			if ef.ProviderDefaults == nil {
				return fmt.Errorf("error: machine specifications for architecture (%s) does not exist nor do provider defaults", d.Architecture)
			}
			val = MachineSpecs{
				// set this to the empty string instead of the data.Architecture
				// so that empty Architectures can be assumed to use provider defaults
				Architecture: "",
				MinWatts:     ef.MinWatts,
				MaxWatts:     ef.MaxWatts,
				// TODO: Need a default/fallback value for GB per chip emissions for AWS and Azure
				GBPerChip: 0,
			}
		}
		d.MachineSpecs = val
		ef.Embodied[d.MachineType] = d
	}
	return nil
}

// readYamlData reads a yaml file and returns a slice of bytes
func readYamlData(filePath string, data interface{}) error {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, data)
	if err != nil {
		return err
	}

	return nil
}
