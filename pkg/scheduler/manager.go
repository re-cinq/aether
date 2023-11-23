package scheduler

import (
	"github.com/re-cinq/cloud-carbon/pkg/bus"
	amazon "github.com/re-cinq/cloud-carbon/pkg/providers/aws"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
)

type ScrapingManager struct {
	schedulers []v1.Scheduler
}

func NewScrapingManager(eventBus *bus.EventBus) *ScrapingManager {

	return &ScrapingManager{
		schedulers: []v1.Scheduler{
			amazon.NewScheduler(eventBus),
		},
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
