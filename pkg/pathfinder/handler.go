package pathfinder

import (
	"github.com/re-cinq/cloud-carbon/pkg/bus"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	"k8s.io/klog/v2"
)

type PahfinderEventHandler struct {
	eventBus *bus.EventBus
}

func NewPahfinderEventHandler(eventBus *bus.EventBus) *PahfinderEventHandler {
	return &PahfinderEventHandler{
		eventBus: eventBus,
	}
}

func (c *PahfinderEventHandler) Apply(event bus.Event) {

	// Make sure we got the right event
	if _, ok := event.(*v1.EmissionsCalculated); ok {

		// TODO: update the prometheus registry

		return
	}

	klog.Errorf("PahfinderEventHandler got an unknown event: %+v", event)

}
