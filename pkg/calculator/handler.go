package calculator

import (
	"github.com/re-cinq/cloud-carbon/pkg/bus"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	"k8s.io/klog/v2"
)

type EmissionCalculator struct {
	eventBus *bus.EventBus
}

func NewEmissionCalculator(eventBus *bus.EventBus) *EmissionCalculator {
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
