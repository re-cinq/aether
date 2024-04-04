package plugin

import (
	"context"

	"github.com/re-cinq/aether/pkg/bus"
	"github.com/re-cinq/aether/pkg/log"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

type PluginHandler struct {
	exporters *ExportPluginSystem
}

func NewHandler(ctx context.Context, e *ExportPluginSystem) *PluginHandler {
	return &PluginHandler{exporters: e}
}

func (p *PluginHandler) Handle(ctx context.Context, e *bus.Event) {
	switch e.Type {
	case v1.EmissionsCalculatedEvent:
		p.SendToExporters(ctx, e)
	default:
		return
	}
}

func (p *PluginHandler) Stop(ctx context.Context) {}

func (p *PluginHandler) SendToExporters(ctx context.Context, e *bus.Event) {
	logger := log.FromContext(ctx)
	instance, ok := e.Data.(v1.Instance)
	if !ok {
		// wrong data on event
		return
	}

	for i := range p.exporters.Plugins {
		err := p.exporters.Plugins[i].Send(&instance)
		if err != nil {
			logger.Error("exporting instance to plugin failed", "error", err)
		}
	}
}
