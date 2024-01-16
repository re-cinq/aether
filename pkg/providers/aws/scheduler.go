package amazon

import (
	"context"
	"time"

	"github.com/re-cinq/cloud-carbon/pkg/config"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	bus "github.com/re-cinq/go-bus"
	"k8s.io/klog/v2"
)

type scheduler struct {

	// Ticker
	ticker *time.Ticker

	// Signal we are done and shutting down
	done chan bool

	// Regions to scrape
	regions []string

	// Event bus
	eventBus bus.Bus

	// AWS Client
	client *Client
}

// Return the scheduler interface
func NewScheduler(ctx context.Context, eventBus bus.Bus) []v1.Scheduler {
	// Load the config
	cfg, exists := config.AppConfig().Providers[provider]
	if !exists {
		return nil
	}

	// Schedulers for each account
	var schedulers []v1.Scheduler

	for index := range cfg.Accounts {
		account := cfg.Accounts[index]

		// Init the AWS client
		client, err := NewClient(ctx, &account, nil)
		if err != nil {
			klog.Errorf("failed to Initialize AWS provider %s", err)
			return nil
		}

		// Init the ticket
		ticker := time.NewTicker(config.AppConfig().ProvidersConfig.Interval)

		// Get the list of regions
		regions := account.Regions

		// Build the initial cache of instances
		for _, region := range regions {
			err := client.ec2Client.Refresh(ctx, client.cache, region)
			if err != nil {
				klog.Errorf("error refreshing EC2 cache at region: %s: %s", region, err)
				continue
			}
		}

		// Build the scheduler
		s := scheduler{
			ticker:   ticker,
			done:     make(chan bool),
			regions:  regions,
			eventBus: eventBus,
			client:   client,
		}

		// Append the scheduler
		schedulers = append(schedulers, &s)
	}

	return schedulers
}

func (s *scheduler) process(ctx context.Context) {
	if len(s.regions) == 0 {
		klog.Error("no AWS regions defined in the config")
		return
	}

	interval := config.AppConfig().Interval

	for _, region := range s.regions {
		// refresh instance cache
		if err := s.client.ec2Client.Refresh(ctx, s.client.cache, region); err != nil {
			klog.Errorf("error refreshing EC2 instances: %s", err)
			return
		}

		instances, err := s.client.cloudWatchClient.GetEC2Metrics(
			s.client.cache,
			region,
			interval,
		)
		if err != nil {
			klog.Errorf("error getting EC2 Metrics with cloudwatch: %s", err)
			return
		}

		for i := range instances {
			// Publish the metrics
			s.eventBus.Publish(v1.MetricsCollected{
				Instance: instances[i],
			})
		}
	}
}

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

	klog.Info("started AWS scheduling")

	// Do the first call
	s.process(ctx)
}

func (s *scheduler) Cancel() {
	// We are done
	s.done <- true

	// Stop the ticker
	s.ticker.Stop()
}
