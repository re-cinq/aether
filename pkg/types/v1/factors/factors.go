package v1

import (
	"fmt"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"

	"github.com/go-yaml/yaml"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

const (
	emissionDataRepoURL = "https://github.com/re-cinq/emissions-data/"
	repoPath            = "/tmp/emissions-data/"
)

var DataPath string = fmt.Sprintf("%s/data/v1", repoPath)

// Emission data is currently stored as files in our emissions-data repo.
// Each file is named "{provider}-{emissionFactor}" where emissionFactor
// may be default, embodied, grid, and use.

// ProviderEmissions reads in emission data for a specified
// provider and stores them into the emissionFactors struct for calulating
func ProviderEmissions(provider v1.Provider, dataPath string) (*EmissionFactors, error) {
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

// CloneAndUpdateFactorsData wraps the CloneAndUpdateRepo
// function with private variables passed.
func CloneAndUpdateFactorsData() error {
	return CloneAndUpdateRepo(repoPath, emissionDataRepoURL)
}

// CloneAndUpdateRepo checks if a local repo exists and is
// up to date with origin. Otherwise, it deletes and clones
// it to the repoPath.
func CloneAndUpdateRepo(repoPath, repoURL string) error {
	// Get repo info if it exists locally
	repo, err := git.PlainOpen(repoPath)
	if err == nil {
		// repo exists, check if its up to date with upstream
		errFetch := repo.Fetch(&git.FetchOptions{Depth: 1})
		if errFetch == git.NoErrAlreadyUpToDate {
			// repo cloned and up to date
			return nil
		}
		// local repo outdated, so remove the directory
		// to be re-cloned later
		if errFetch == nil {
			os.RemoveAll(repoPath)
		} else {
			return err
		}
	}

	// Repo doesn't exist, so clone it
	if err == git.ErrRepositoryNotExists {
		_, errClone := git.PlainClone(repoPath, false, &git.CloneOptions{
			URL:   repoURL,
			Depth: 1,
		})
		return errClone
	}
	return err
}
