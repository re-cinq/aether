package exporter

import (
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	bus "github.com/re-cinq/go-bus"
	"k8s.io/klog/v2"
)

type PrometheusEventHandler struct {
	eventBus bus.Bus
}

func NewPrometheusEventHandler(eventBus bus.Bus) *PrometheusEventHandler {
	return &PrometheusEventHandler{
		eventBus: eventBus,
	}
}

func (c *PrometheusEventHandler) Apply(event bus.Event) {

	// Make sure we got the right event
	if _, ok := event.(*v1.EmissionsCalculated); ok {

		// TODO: update the prometheus registry

		return
	}

	klog.Errorf("PrometheusEventHandler got an unknown event: %+v", event)

}
