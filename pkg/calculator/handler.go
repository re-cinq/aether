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
	return &EmissionCalculator{
		eventBus: eventBus,
	}
}

func (ec *EmissionCalculator) Apply(event bus.Event) {
	// Make sure we got the right event
	if metricsCollected, ok := event.(v1.MetricsCollected); ok {
		instance := metricsCollected.Instance
		cfg := config.AppConfig().ProvidersConfig

		emFactors, err := factors.GetEmissionFactors(instance.Provider(), cfg.FactorsDataPath)
		if err != nil {
			klog.Errorf("error getting emission factors: %+v", err)
			return
		}

		specs, ok := emFactors.Embodied[instance.Kind()]
		if !ok {
			klog.Errorf("error finding instance: %s in factor data", instance.Name())
			return
		}

		mCPU, ok := instance.Metrics()[v1.CPU.String()]
		if !ok {
			klog.Errorf("error instance metrics for CPU don't exist")
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
			cores:         mCPU.UnitAmount(),
			usage:         mCPU.Usage(),
			pue:           emFactors.AveragePUE,
			gridCO2e:      gridCO2e,
		}

		instance.SetOperationalEmissions(
			v1.NewResourceEmission(
				c.operationalEmissions(cfg.Interval),
				v1.GCO2eqkWh,
			),
		)

		instance.SetEmbodiedEmissions(
			v1.NewResourceEmission(
				c.embodiedEmissions(cfg.Interval),
				v1.GCO2eqkWh,
			),
		)

		ec.eventBus.Publish(v1.EmissionsCalculated{
			Instance: instance,
		})

		for _, metric := range instance.Metrics() {
			klog.Infof("Collected metric: %s %s %s %s | %s", instance.Service(), instance.Region(), instance.Name(), instance.Kind(), metric.String())
		}

		return
	}

	klog.Errorf("EmissionCalculator got an unknown event: %+v", event)
}
