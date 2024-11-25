package gcp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/re-cinq/aether/pkg/config"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

// Source is a configured google source that adheres to Aethers source
// interface
type Source struct {
	// Google Client
	*Client

	// The Project that Client is setup for
	Project *string

	// Teardown functionality
	Shutdown func()
}

// Sources instantiates a slice of instances of the Google Sources configured
// for use, one per project configured
func Sources(ctx context.Context, cfg *config.Provider) []v1.Source {
	var sources []v1.Source

	// we instantiate a source per project
	for index := range cfg.Accounts {
		account := cfg.Accounts[index]

		c, shutdown, err := New(ctx, &account)
		if err != nil {
			return nil
		}

		sources = append(sources, &Source{
			Project:  &account.Project,
			Client:   c,
			Shutdown: shutdown,
		})
	}

	return sources
}

// Fetch returns a slice of instances, this is to adhere to the sources
// interface
func (s *Source) Fetch(ctx context.Context) (map[string]*v1.Instance, error) {
	if s.Project == nil {
		return nil, errors.New("no project set")
	}

	s.Client.Refresh(ctx, *s.Project)

	interval := config.AppConfig().Interval
	if interval < 5*time.Minute {
		return nil, fmt.Errorf("error interval for GCP needs to be atleast 5m. It is: %+v", interval)
	}

	err := s.Client.GetMetricsForInstances(ctx, *s.Project, interval.String())
	if err != nil {
		return nil, fmt.Errorf("failed getting instance metrics: %v", err)
	}

	for k, instance := range s.Client.instancesMap {
		// remove terminated instances as we
		// shouldnt use them anymore
		if instance.Status == v1.InstanceTerminated {
			delete(s.Client.instancesMap, k)
		}
	}

	return s.Client.instancesMap, nil
}

// Stop is used to gracefully shutdown a source
func (s *Source) Stop(ctx context.Context) error {
	s.Shutdown()
	return nil
}
