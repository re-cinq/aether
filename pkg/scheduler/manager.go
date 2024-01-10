package scheduler

import (
	"context"

	amazon "github.com/re-cinq/cloud-carbon/pkg/providers/aws"
	"github.com/re-cinq/cloud-carbon/pkg/providers/gcp"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	bus "github.com/re-cinq/go-bus"
)

type ScrapingManager struct {
	schedulers []v1.Scheduler
}

func NewScrapingManager(ctx context.Context, eventBus bus.Bus) *ScrapingManager {
	var schedulers []v1.Scheduler

	// Add AWS
	schedulers = append(schedulers, amazon.NewScheduler(eventBus)...)

	// Add GCP
	schedulers = append(schedulers, gcp.NewScheduler(ctx, eventBus)...)

	return &ScrapingManager{
		schedulers: schedulers,
	}
}

func (m ScrapingManager) Start(ctx context.Context) {
	for _, scheduler := range m.schedulers {
		if scheduler != nil {
			scheduler.Schedule(ctx)
		}
	}
}

func (m ScrapingManager) Stop() {
	for _, scheduler := range m.schedulers {
		if scheduler != nil {
			scheduler.Cancel()
		}
	}
}
