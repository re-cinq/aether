package calculator

import (
	"fmt"
	"log/slog"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/re-cinq/cloud-carbon/pkg/config"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	factors "github.com/re-cinq/cloud-carbon/pkg/types/v1/factors"
	data "github.com/re-cinq/emissions-data/pkg/types/v2"
	bus "github.com/re-cinq/go-bus"
)

var awsInstances map[string]data.Instance

// AWS, GCP and Azure have increased their server lifespan to 6 years (2024)
// https://sustainability.aboutamazon.com/products-services/the-cloud?energyType=true
// https://www.theregister.com/2024/01/31/alphabet_q4_2023/
// https://www.theregister.com/2022/08/02/microsoft_server_life_extension/
const serverLifespan = 6

type EmissionCalculator struct {
	eventBus bus.Bus
}

func NewEmissionCalculator(eventBus bus.Bus) *EmissionCalculator {
	err := factors.CloneAndUpdateFactorsData()
	if err != nil {
		slog.Error("error with emissions repo", "error", err)
		return nil
	}

	awsInstances, err = getProviderEC2EmissionFactors(v1.Aws)
	if err != nil {
		slog.Error("unable to get v2 Emission Factors, falling back to v1", "error", err)
	}

	return &EmissionCalculator{
		eventBus: eventBus,
	}
}

func getProviderEC2EmissionFactors(provider v1.Provider) (map[string]data.Instance, error) {
	yamlURL := fmt.Sprintf("https://raw.githubusercontent.com/re-cinq/emissions-data/main/data/v2/%s-instances.yaml", provider)
	resp, err := http.Get(yamlURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder := yaml.NewDecoder(resp.Body)
	err = decoder.Decode(&awsInstances)
	if err != nil {
		return nil, err
	}

	return awsInstances, nil
}

func (ec *EmissionCalculator) Apply(event bus.Event) {
	interval := config.AppConfig().ProvidersConfig.Interval

	// Make sure we got the right event
	metricsCollected, ok := event.(v1.MetricsCollected)
	if !ok {
		slog.Error("EmissionCalculator got an unknown event", "event", event)
		return
	}
	eventInstance := metricsCollected.Instance

	// Gets PUE, grid data, and machine specs
	emFactors, err := factors.GetProviderEmissionFactors(
		eventInstance.Provider(),
		factors.DataPath,
	)
	if err != nil {
		slog.Error("error getting emission factors", "error", err)
		return
	}

	gridCO2eTons, ok := emFactors.Coefficient[eventInstance.Region()]
	if !ok {
		slog.Error("region does not exist in factors for provider", "region", eventInstance.Region(), "provider", "gcp")
		return
	}
	// TODO: hotfix until updated in emissions data
	// convert gridCO2e from metric tonnes to grams
	gridCO2e := gridCO2eTons * (1000 * 1000)

	params := parameters{
		gridCO2e: gridCO2e,
		pue:      emFactors.AveragePUE,
	}

	specs, ok := emFactors.Embodied[eventInstance.Kind()]
	if !ok {
		slog.Error("failed finding instance in factor data", "instance", eventInstance.Name(), "kind", eventInstance.Kind())
		return
	}

	if d, ok := awsInstances[eventInstance.Kind()]; ok {
		params.wattage = d.PkgWatt
		params.vCPU = float64(d.VCPU)
		params.embodiedFactor = d.EmbodiedHourlyGCO2e
	} else {
		params.wattage = []data.Wattage{
			{
				Percentage: 0,
				Wattage:    specs.MinWatts,
			},
			{
				Percentage: 100,
				Wattage:    specs.MaxWatts,
			},
		}
		params.embodiedFactor = hourlyEmbodiedEmissions(&specs)
	}

	// calculate and set the operational emissions for each
	// metric type (CPU, Memory, Storage, and networking)
	metrics := eventInstance.Metrics()
	for _, v := range metrics {
		params.metric = &v
		opEm, err := operationalEmissions(interval, &params)
		if err != nil {
			slog.Error("failed calculating operational emissions", "type", v.Name(), "error", err)
			continue
		}
		params.metric.SetEmissions(v1.NewResourceEmission(opEm, v1.GCO2eqkWh))
		// update the instance metrics
		metrics.Upsert(params.metric)
	}

	eventInstance.SetEmbodiedEmissions(
		v1.NewResourceEmission(
			embodiedEmissions(interval, params.embodiedFactor),
			v1.GCO2eqkWh,
		),
	)

	ec.eventBus.Publish(v1.EmissionsCalculated{
		Instance: eventInstance,
	})

	eventInstance.PrintPretty()
}

func hourlyEmbodiedEmissions(e *factors.Embodied) float64 {
	// we fall back on the specs from the previous dataset
	// and convert it into a hourly factor
	// this is based on CCF's calculation:
	//
	// M = TE * (TR/EL) * (RR/TR)
	//
	// TE = Total Embodied Emissions
	// TR = Time Reserved (in years)
	// EL = Expected Lifespan
	// RR = Resources Reserved
	// TR = Total Resources, the total number of resources available.
	return e.TotalEmbodiedKiloWattCO2e *
		// 1 hour normalized to a year
		((1.0 / 24.0 / 365.0) / serverLifespan) *
		// amount of vCPUS for instance versus total vCPUS for platform
		(e.VCPU / e.TotalVCPU)
}
