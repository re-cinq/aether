package pathfinder

import (
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	bus "github.com/re-cinq/go-bus"
	"k8s.io/klog/v2"
)

type PathfinderEventHandler struct {
	eventBus bus.Bus
}

func NewPathfinderEventHandler(eventBus bus.Bus) *PathfinderEventHandler {
	return &PathfinderEventHandler{
		eventBus: eventBus,
	}
}

func (c *PathfinderEventHandler) Apply(event bus.Event) {
	// Make sure we got the right event
	if _, ok := event.(*v1.EmissionsCalculated); ok {
		// TODO: update the prometheus registry

		return
	}

	klog.Errorf("PathfinderEventHandler got an unknown event: %+v", event)
}
