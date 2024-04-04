package amazon

import (
	"context"
	"errors"
	"fmt"
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
			err := c.Refresh(ctx, region)
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
				err := s.scrape(ctx)
				if err != nil {
					s.logger.Error("scraping error", "error", err)
				}
			}
		}
	}()

	// do the first scrape
	err := s.scrape(ctx)
	if err != nil {
		s.logger.Error("scraping error", "error", err)
	}
}

func (s *Scraper) scrape(ctx context.Context) error {
	if len(s.regions) == 0 {
		return errors.New("no AWS regions defined in the config")
	}

	interval := config.AppConfig().Interval

	for _, region := range s.regions {
		// refresh instance cache
		if err := s.Client.Refresh(ctx, region); err != nil {
			return err
		}

		err := s.Client.GetEC2Metrics(
			ctx,
			region,
			interval,
		)
		if err != nil {
			return err
		}
	}

	var errs error
	for _, instance := range s.Client.instancesMap {
		err := s.Bus.Publish(&bus.Event{
			Type: v1.MetricsCollectedEvent,
			Data: *instance,
		})
		if err != nil {
			errs = errors.Join(
				errs,
				fmt.Errorf("failed to publish for instance %s: %v", instance.Name, err),
			)
		}
	}
	return errs
}

func (s *Scraper) Stop(ctx context.Context) {
	s.Done <- true

	s.ticker.Stop()
}
