package calculator

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"

	"github.com/re-cinq/cloud-carbon/pkg/config"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	factors "github.com/re-cinq/cloud-carbon/pkg/types/v1/factors"
	data "github.com/re-cinq/emissions-data/pkg/types/v2"
	bus "github.com/re-cinq/go-bus"
)

var awsInstances map[string]data.Instance

type EmissionCalculator struct {
	eventBus bus.Bus
}

func NewEmissionCalculator(eventBus bus.Bus) *EmissionCalculator {
	err := factors.CloneAndUpdateFactorsData()
	if err != nil {
		klog.Errorf("error with emissions repo: %+v", err)
		return nil
	}

	awsInstances, err = getProviderEC2EmissionFactors(v1.Aws)
	if err != nil {
		klog.Error("unable to get v2 Emission Factors, falling back to v1: ", err)
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
		klog.Errorf("EmissionCalculator got an unknown event: %+v", event)
		return
	}
	eventInstance := metricsCollected.Instance

	// Gets PUE, grid data, and machine specs
	emFactors, err := factors.GetProviderEmissionFactors(
		eventInstance.Provider(),
		factors.DataPath,
	)
	if err != nil {
		klog.Errorf("error getting emission factors: %+v", err)
		return
	}

	gridCO2e, ok := emFactors.Coefficient[eventInstance.Region()]
	if !ok {
		klog.Errorf("error region: %s does not exist in factors for %s", eventInstance.Region(), "gcp")
		return
	}

	params := parameters{
		gridCO2e: gridCO2e,
		pue:      emFactors.AveragePUE,
	}

	specs, ok := emFactors.Embodied[eventInstance.Kind()]
	if !ok {
		klog.Errorf("error finding instance: %s kind: %s in factor data", eventInstance.Name(), eventInstance.Kind())
		return
	}

	if d, ok := awsInstances[eventInstance.Kind()]; ok {
		params.wattage = d.PkgWatt
		params.vCPU = float64(d.VCPU)
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
	}

	// calculate and set the operational emissions for each
	// metric type (CPU, Memory, Storage, and networking)
	metrics := eventInstance.Metrics()
	for _, v := range metrics {
		params.metric = v
		opEm, err := operationalEmissions(interval, &params)
		if err != nil {
			klog.Errorf("error calculating %s operational emissions: %+v", v.Name(), err)
			continue
		}
		params.metric.SetEmissions(v1.NewResourceEmission(opEm, v1.GCO2eqkWh))
		// update the instance metrics
		metrics.Upsert(&params.metric)
	}

	eventInstance.SetEmbodiedEmissions(
		v1.NewResourceEmission(
			embodiedEmissions(interval, specs.TotalEmbodiedKiloWattCO2e),
			v1.GCO2eqkWh,
		),
	)

	ec.eventBus.Publish(v1.EmissionsCalculated{
		Instance: eventInstance,
	})

	eventInstance.PrintPretty()
}
