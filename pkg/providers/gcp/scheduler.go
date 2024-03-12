package gcp

import (
	"context"
	"log/slog"
	"time"

	"github.com/re-cinq/cloud-carbon/pkg/config"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	bus "github.com/re-cinq/go-bus"
)

type scheduler struct {
	// Ticker
	ticker *time.Ticker

	// Signal we are done and shutting down
	done chan bool

	// Event bus
	eventBus bus.Bus

	// Project ID
	project string

	// GCP handler
	gcp *GCP

	// Shutdown function
	shutdown func()
}

// NewScheduler returns the a list of schedulers that
// conform to the scheduler interface
func NewScheduler(ctx context.Context, eventBus bus.Bus) []v1.Scheduler {
	// Load the GCP config
	cfg, exists := config.AppConfig().Providers[provider]

	// If the provider is not configured - skip its initialization
	if !exists {
		return nil
	}

	var schedulers []v1.Scheduler

	// Create a scheduler for each GCP project
	for index := range cfg.Accounts {
		account := cfg.Accounts[index]

		ticker := time.NewTicker(config.AppConfig().ProvidersConfig.Interval)

		// Init the GCP Client
		gcp, shutdown, err := New(ctx, &account)
		if err != nil {
			slog.Error("failed to Initialize GCP provider", "error", err)
			return nil
		}

		// loads all instances into cache
		gcp.Refresh(ctx, account.Project)

		schedulers = append(schedulers, &scheduler{
			ticker:   ticker,
			done:     make(chan bool),
			project:  account.Project,
			eventBus: eventBus,
			gcp:      gcp,
			shutdown: shutdown,
		})
	}

	return schedulers
}

// process is the logic for the scheduler
// to run at certain intervals
func (s *scheduler) process(ctx context.Context) {
	if s.project == "" {
		slog.Error("no GCP project defined in the config")
		return
	}

	s.gcp.Refresh(ctx, s.project)

	instances, err := s.gcp.GetMetricsForInstances(ctx, s.project, "5m")

	if err != nil {
		slog.Error("failed to scrape instance metrics", "error", err)
		return
	}

	for i := range instances {
		instance := instances[i]
		s.eventBus.Publish(v1.MetricsCollected{
			Instance: instance,
		})
	}
}

// Schedule setups a schedule and runs it according to a time intervals
// only stops when receives a done signal
func (s *scheduler) Schedule(ctx context.Context) {
	go func() {
		for {
			select {
			case <-s.done:
				return
			case <-s.ticker.C:
				s.process(ctx)
			}
		}
	}()

	slog.Info("started GCP scheduling")

	// Do the first call
	s.process(ctx)
}

// Cancel stops the scheduler
func (s *scheduler) Cancel() {
	// We are done
	s.done <- true

	// Stop the ticker
	s.ticker.Stop()

	// Shutdown the GCP handler
	s.shutdown()
}
