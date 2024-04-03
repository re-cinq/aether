package amazon

import (
	"context"
	"log/slog"
	"time"

	"github.com/re-cinq/aether/pkg/bus"
	"github.com/re-cinq/aether/pkg/config"
	"github.com/re-cinq/aether/pkg/log"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

type Scraper struct {
	// Ticker
	ticker *time.Ticker
	Done   chan bool

	Client *Client

	// Regions to scrape
	regions []string

	Bus *bus.Bus

	logger *slog.Logger
}

func SetupScrapers(ctx context.Context, b *bus.Bus) []v1.Scraper {
	cfg, exists := config.AppConfig().Providers[provider]
	logger := log.FromContext(ctx)

	// If the provider is not configured - skip its initialization
	if !exists {
		return nil
	}

	var scrapers []v1.Scraper

	for index := range cfg.Accounts {
		account := cfg.Accounts[index]

		ticker := time.NewTicker(config.AppConfig().ProvidersConfig.Interval)

		c, err := New(ctx, &account, nil)
		if err != nil {
			return nil
		}

		// Get the list of regions
		regions := account.Regions

		// Build the initial cache of instances
		for _, region := range regions {
			err := c.ec2Client.Refresh(ctx, c.cache, region)
			if err != nil {
				logger.Error("error refreshing cache for region", "region", region, "error", err)
				continue
			}
		}

		scrapers = append(scrapers, &Scraper{
			ticker:  ticker,
			Done:    make(chan bool),
			regions: regions,
			Bus:     b,
			Client:  c,
			logger:  logger,
		})
	}

	return scrapers
}

func (s *Scraper) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-s.Done:
				return
			case <-s.ticker.C:
				s.scrape(ctx)
			}
		}
	}()

	// do the first scrape
	s.scrape(ctx)
}

func (s *Scraper) scrape(ctx context.Context) {
	if len(s.regions) == 0 {
		s.logger.Error("no AWS regions defined in the config")
		return
	}

	interval := config.AppConfig().Interval

	for _, region := range s.regions {
		// refresh instance cache
		if err := s.Client.ec2Client.Refresh(ctx, s.Client.cache, region); err != nil {
			s.logger.Error("error refreshing EC2 instances", "error", err)
			return
		}

		instances, err := s.Client.cloudWatchClient.GetEC2Metrics(
			ctx,
			s.Client.cache,
			region,
			interval,
		)
		if err != nil {
			s.logger.Error("error getting EC2 Metrics with cloudwatch", "error", err)
			return
		}

		for i := range instances {
			// Publish the metrics
			if err := s.Bus.Publish(&bus.Event{
				Type: v1.MetricsCollectedEvent,
				Data: instances[i],
			}); err != nil {
				s.logger.Error("failed publishing instance", "error", err, "instance", instances[i].Name)
			}
		}
	}
}

func (s *Scraper) Stop(ctx context.Context) {
	s.Done <- true

	s.ticker.Stop()
}
