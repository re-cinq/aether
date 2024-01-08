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
	_, ok := event.(v1.EmissionsCalculated)
	if !ok {
		klog.Errorf("PathfinderEventHandler got an unknown event: %+v", event)
		return
	}
	// TODO: send data to the pathfinder api
}
