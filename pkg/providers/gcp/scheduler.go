package gcp

import (
	"context"
	"time"

	"github.com/re-cinq/cloud-carbon/pkg/config"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	bus "github.com/re-cinq/go-bus"
	"k8s.io/klog/v2"
)

type gcpScheduler struct {

	// Ticker
	ticker *time.Ticker

	// Signal we are done and shutting down
	done chan bool

	// Event bus
	eventBus bus.Bus

	// Project ID
	project string

	// GCP Metrics client
	gcpClient *GCP

	// GCE client
	gce *gceClient

	// Account config
	account config.Account

	// Shutdown function
	shutdown func()
}

// Return the scheduler interface
func NewScheduler(eventBus bus.Bus) []v1.Scheduler {
	// Load the config
	gcpConfig, exists := config.AppConfig().Providers[gcpProvider]

	// If the provider is not configured - skip its initialization
	if !exists {
		return nil
	}

	// Schedulers for each account
	var schedulers []v1.Scheduler

	for index := range gcpConfig.Accounts {
		account := gcpConfig.Accounts[index]

		// Init the ticket
		ticker := time.NewTicker(config.AppConfig().ProvidersConfig.Interval)

		// Init the GCE client
		gce := newGCECLient(&account)
		if gce == nil {
			klog.Error("failed to Initialize GCP provider")
			return nil
		}

		// Init the GCP metrics Client
		// This doesn't return metrics, it returns a new GCP client
		gcpClient, shutdown, err := New(&account, gce.cache)
		if err != nil {
			klog.Errorf("failed to Initialize GCP provider %s", err)
			return nil
		}

		// List all the instances
		gce.Refresh(account.Project)

		accountScheduler := gcpScheduler{
			ticker:    ticker,
			done:      make(chan bool),
			project:   account.Project,
			account:   account,
			eventBus:  eventBus,
			gcpClient: gcpClient,
			gce:       gce,
			shutdown:  shutdown,
		}

		schedulers = append(schedulers, &accountScheduler)
	}

	return schedulers
}

func (s *gcpScheduler) process() {
	if s.project == "" {
		klog.Error("no GCP project defined in the config")
		return
	}

	instances, err := s.gcpClient.GetMetricsForInstances(context.TODO(), "5m")

	if err != nil {
		klog.Errorf("failed to scrape instance metrics %s", err)
		return
	}

	for _, instance := range instances {
		// // Publish the metrics
		s.eventBus.Publish(v1.MetricsCollected{
			Instance: instance,
		})
	}
}

func (s *gcpScheduler) Schedule() {
	go func() {
		for {
			select {
			case <-s.done:
				return
			case <-s.ticker.C:
				s.process()
			}
		}
	}()

	klog.Info("started GCP scheduling")

	// Do the first call
	s.process()
}

func (s *gcpScheduler) Cancel() {
	// We are done
	s.done <- true

	// Stop the ticker
	s.ticker.Stop()

	// Shutdown the GCP client
	s.shutdown()

	// Close the GCE client
	s.gce.Close()
}
