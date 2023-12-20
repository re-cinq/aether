package calculator

import (
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
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

func (c *EmissionCalculator) Apply(event bus.Event) {
	// Make sure we got the right event
	if metricsCollected, ok := event.(v1.MetricsCollected); ok {
		// TODO: remove this, which is only for debugging purposes
		instance := metricsCollected.Instance

		// TODO: do the calculation

		// TODO: emit the calculation event in the bus

		for _, metric := range instance.Metrics() {
			klog.Infof("Collected metric: %s %s %s %s | %s", instance.Service(), instance.Region(), instance.Name(), instance.Kind(), metric.String())
		}

		return
	}

	klog.Errorf("EmissionCalculator got an unknown event: %+v", event)
}
