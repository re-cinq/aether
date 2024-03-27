package scraper

import (
	"context"

	"github.com/re-cinq/aether/pkg/bus"
	amazon "github.com/re-cinq/aether/pkg/providers/aws"
	"github.com/re-cinq/aether/pkg/providers/gcp"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

// ScraperManager used to handle the various scrapers
type ScrapingManager struct {
	scrapers []v1.Scraper
}

// NewManager returns a configured instance of a ScraperManager
func NewManager(ctx context.Context, b *bus.Bus) *ScrapingManager {
	var scrapers []v1.Scraper

	// Add aws
	scrapers = append(scrapers, amazon.SetupScrapers(ctx, b)...)
	// Add GCP
	scrapers = append(scrapers, gcp.SetupScrapers(ctx, b)...)

	return &ScrapingManager{scrapers}
}

// Start iterates through each scraper and starts them
// NOTE: a scraper should not have a blocking call to start
func (m ScrapingManager) Start(ctx context.Context) {
	for _, scraper := range m.scrapers {
		if scraper != nil {
			scraper.Start(ctx)
		}
	}
}

// Stop iterates through scrapers and stops them
func (m ScrapingManager) Stop(ctx context.Context) {
	for _, scraper := range m.scrapers {
		if scraper != nil {
			scraper.Stop(ctx)
		}
	}
}
