package pathfinder

import (
	"context"
	"log/slog"

	"github.com/re-cinq/cloud-carbon/pkg/log"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	bus "github.com/re-cinq/go-bus"
)

type PathfinderEventHandler struct {
	eventbus bus.Bus
	logger   *slog.Logger
}

func NewPathfinderEventHandler(ctx context.Context, eventbus bus.Bus) *PathfinderEventHandler {
	logger := log.FromContext(ctx)

	return &PathfinderEventHandler{eventbus, logger}
}

func (p *PathfinderEventHandler) Apply(event bus.Event) {

	// Make sure we got the right event
	_, ok := event.(v1.EmissionsCalculated)
	if !ok {
		p.logger.Error("PathfinderEventHandler got an unknown event", "event", event)
		return
	}
	// TODO: send data to the pathfinder api
}
