package amazon

import (
	"context"
	"errors"
	"fmt"

	"github.com/re-cinq/aether/pkg/config"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

// Source is a configured Amazon source that adheres to Aethers source
// interface
type Source struct {
	// Amazon Client
	*Client

	// AWS doesn't use project, but instead is
	// separated by region
	Region string
}

// Sources instantiates a slice of instances of the Amazon Sources configured
// for use, one per project configured
func Sources(ctx context.Context, cfg *config.Provider) []v1.Source {
	var sources []v1.Source

	// we instantiate a source per project
	for index := range cfg.Accounts {
		account := cfg.Accounts[index]

		c, err := New(ctx, &account, nil)
		if err != nil {
			return nil
		}

		for _, region := range account.Regions {
			sources = append(sources, &Source{
				Region: region,
				Client: c,
			})
		}
	}

	return sources
}

// Fetch returns a slice of instances, this is to adhere to the sources
// interface
func (s *Source) Fetch(ctx context.Context) ([]*v1.Instance, error) {
	if s.Region == "" {
		return nil, errors.New("no region set")
	}

	err := s.Client.Refresh(ctx, s.Region)
	if err != nil {
		return nil, err
	}

	interval := config.AppConfig().Interval

	err = s.Client.GetEC2Metrics(ctx, s.Region, interval)
	if err != nil {
		return nil, fmt.Errorf("failed getting instance metrics: %v", err)
	}

	var instances []*v1.Instance
	for _, instance := range s.Client.instancesMap {
		instances = append(instances, instance)
	}

	return instances, nil
}

// Stop is used to gracefully shutdown a source
func (s *Source) Stop(ctx context.Context) error {
	return nil
}
