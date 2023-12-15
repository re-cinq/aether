package scheduler

import (
	amazon "github.com/re-cinq/cloud-carbon/pkg/providers/aws"
	"github.com/re-cinq/cloud-carbon/pkg/providers/gcp"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	bus "github.com/re-cinq/go-bus"
)

type ScrapingManager struct {
	schedulers []v1.Scheduler
}

func NewScrapingManager(eventBus bus.Bus) *ScrapingManager {

	var schedulers []v1.Scheduler

	// Add AWS
	schedulers = append(schedulers, amazon.NewScheduler(eventBus)...)

	// Add GCP
	schedulers = append(schedulers, gcp.NewScheduler(eventBus)...)

	return &ScrapingManager{
		schedulers: schedulers,
	}

}

func (m ScrapingManager) Start() {

	for _, scheduler := range m.schedulers {
		if scheduler != nil {
			scheduler.Schedule()
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
