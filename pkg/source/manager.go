package source

import (
	"context"
	"sync"
	"time"

	"github.com/re-cinq/aether/pkg/bus"
	"github.com/re-cinq/aether/pkg/config"
	"github.com/re-cinq/aether/pkg/log"
	"github.com/re-cinq/aether/pkg/plugin"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

type Manager struct {
	ticker *time.Ticker

	bus *bus.Bus

	Sources []v1.Source

	plugin *plugin.SourcePluginSystem
}

type option func(m *Manager)

func WithPlugins(s *plugin.SourcePluginSystem) option {
	return func(m *Manager) {
		m.plugin = s
	}
}

func WithBus(b *bus.Bus) option {
	return func(m *Manager) {
		m.bus = b
	}
}

func New(ctx context.Context, opts ...option) *Manager {
	m := &Manager{
		ticker: time.NewTicker(config.AppConfig().ProvidersConfig.Interval),
	}

	for _, o := range opts {
		o(m)
	}

	// load buil in sources
	m.Sources = BuiltInSources(ctx)

	// load plugin sources
	for _, p := range m.plugin.Plugins {
		m.Sources = append(m.Sources, p.Source)
	}

	return m
}

// Start is used to start the processing of the manager
func (m *Manager) Start(ctx context.Context) {
	go func() {
		for {
			select {
			// when context is canceled we will
			// stop processing
			case <-ctx.Done():
				return
			case <-m.ticker.C:
				m.Fetch(ctx)
			}
		}
	}()

	// run fetch once the first time
	m.Fetch(ctx)
}

// Fetch goes to all sources and fetchs the instance list from them and then
// publishes all those instances on the bus
func (m *Manager) Fetch(ctx context.Context) {
	logger := log.FromContext(ctx)

	var wg sync.WaitGroup
	for i := range m.Sources {
		wg.Add(1)

		go func(source v1.Source) {
			instances, err := source.Fetch(ctx)
			if err != nil {
				logger.Error("failed fetching instances", "error", err)
				wg.Done()
				return
			}
			logger.Debug("publishing instances", "instance count", len(instances))
			err = m.publishInstances(instances)
			if err != nil {
				logger.Error("failed publishing instances", "error", err)
			}
			wg.Done()
		}(m.Sources[i])
	}

	wg.Wait()
}

// publishInstances is a helper that publishes each instance in a slice on the
// bus under the MetricsCollectedEvent
func (m *Manager) publishInstances(instances []*v1.Instance) error {
	for i := range instances {
		err := m.bus.Publish(&bus.Event{
			Type: v1.MetricsCollectedEvent,
			Data: *instances[i],
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop is used to graceful shut down the manager and by extension all the
// sources
func (m *Manager) Stop(ctx context.Context) {
	logger := log.FromContext(ctx)

	m.ticker.Stop()

	for i := range m.Sources {
		err := m.Sources[i].Stop(ctx)
		if err != nil {
			logger.Error("failed stopping source", "error", err)
		}
	}
}
