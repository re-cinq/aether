package pathfinder

import (
	"context"

	"github.com/re-cinq/aether/pkg/bus"
)

type PathfinderHandler struct{}

func NewHandler() *PathfinderHandler {
	return &PathfinderHandler{}
}

func (p *PathfinderHandler) Handle(ctx context.Context, e *bus.Event) {}

func (p *PathfinderHandler) Stop(ctx context.Context) {}
