package calculator

import (
	"github.com/re-cinq/cloud-carbon/pkg/config"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	factors "github.com/re-cinq/cloud-carbon/pkg/types/v1/factors"
	bus "github.com/re-cinq/go-bus"
	"k8s.io/klog/v2"
)

type EmissionCalculator struct {
	eventBus bus.Bus
}

func NewEmissionCalculator(eventBus bus.Bus) *EmissionCalculator {
	err := factors.CloneAndUpdateFactorsData()
	if err != nil {
		klog.Errorf("error with emissions repo: %+v", err)
		return nil
	}

	return &EmissionCalculator{
		eventBus: eventBus,
	}
}

func (ec *EmissionCalculator) Apply(event bus.Event) {
	// Make sure we got the right event
	metricsCollected, ok := event.(v1.MetricsCollected)
	if !ok {
		klog.Errorf("EmissionCalculator got an unknown event: %+v", event)
		return
	}

	instance := metricsCollected.Instance
	interval := config.AppConfig().ProvidersConfig.Interval

	emFactors, err := factors.GetProviderEmissionFactors(
		instance.Provider(),
		factors.DataPath,
	)
	if err != nil {
		klog.Errorf("error getting emission factors: %+v", err)
		return
	}

	specs, ok := emFactors.Embodied[instance.Kind()]
	if !ok {
		klog.Errorf("error finding instance: %s kind: %s in factor data", instance.Name(), instance.Kind())
		return
	}

	gridCO2e, ok := emFactors.Coefficient[instance.Region()]
	if !ok {
		klog.Errorf("error region: %s does not exist in factors for %s", instance.Region(), "gcp")
		return
	}

	c := calculate{
		minWatts:      specs.MinWatts,
		maxWatts:      specs.MaxWatts,
		totalEmbodied: specs.TotalEmbodiedKiloWattCO2e,
		pue:           emFactors.AveragePUE,
		gridCO2e:      gridCO2e,
	}

	// calculate and set the operational emissions for each
	// metric type (CPU, Memory, Storage, and networking)
	metrics := instance.Metrics()
	for _, v := range metrics {
		// reassign to avoid implicit memory aliasing
		v := v
		em, err := c.operationalEmissions(&v, interval)
		if err != nil {
			klog.Errorf("error calculating %s operational emissions: %+v", v.Name(), err)
			continue
		}
		v.SetEmissions(v1.NewResourceEmission(em, v1.GCO2eqkWh))
		// update the instance metrics
		metrics.Upsert(&v)
	}

	instance.SetEmbodiedEmissions(
		v1.NewResourceEmission(
			c.embodiedEmissions(interval),
			v1.GCO2eqkWh,
		),
	)

	ec.eventBus.Publish(v1.EmissionsCalculated{
		Instance: instance,
	})

	instance.PrintPretty()
}
