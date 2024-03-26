package exporter

import (
	"context"
	"log/slog"

	"github.com/re-cinq/cloud-carbon/pkg/bus"
	"github.com/re-cinq/cloud-carbon/pkg/log"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

// PromHandler is the Event handnelr used to configure prometheus
// metrics
type PromHandler struct {
	Bus    *bus.Bus
	meter  api.Meter
	logger *slog.Logger
}

// NewHandler returns a configured instance of PromHandler
func NewHandler(ctx context.Context, b *bus.Bus) *PromHandler {
	logger := log.FromContext(ctx)

	exporter, err := prometheus.New()
	if err != nil {
		logger.Error("failed setting up prometheus", "error", err)
		return nil
	}

	meter := metric.NewMeterProvider(
		metric.WithReader(exporter),
	).Meter("cloud-carbon")

	return &PromHandler{
		Bus:    b,
		meter:  meter,
		logger: logger,
	}
}

func (p *PromHandler) Stop(ctx context.Context) {}

// Handle handles events passed by the bus, used to adhere to the EventHandler
// intterface
func (p *PromHandler) Handle(ctx context.Context, e *bus.Event) {
	switch e.Type {
	case v1.EmissionsCalculatedEvent:
		p.handleEvent(e)
	default:
		return
	}
}

// handleEvent is the business logic for handleing the v1.EmissionsCalculatedEvent
func (p *PromHandler) handleEvent(e *bus.Event) {
	i, ok := e.Data.(v1.Instance)
	if !ok {
		// wrong data on event
		return
	}

	// setup emissions gauge
	emissions, err := p.meter.Float64ObservableGauge(
		"emissions",
		api.WithDescription("co2eq of various services"),
	)
	if err != nil {
		p.logger.Error("[otel] failed setting up emissions metric")
		return
	}

	// setup embodied emissions gauge
	embodied, err := p.meter.Float64ObservableGauge(
		"embodied",
		api.WithDescription("co2eq of various services"),
	)
	if err != nil {
		p.logger.Error("[otel] failed setting up embodied emissions metric")
		return
	}

	// register embodied emissions metrics for instance
	// NOTE: this will not change based on different types of metrics
	_, err = p.meter.RegisterCallback(
		func(ctx context.Context, o api.Observer) error {
			o.ObserveFloat64(
				embodied,
				i.EmbodiedEmissions.Value,
				api.WithAttributes(
					getAttributesFromInstance(&i)...,
				))

			return nil
		}, embodied)
	if err != nil {
		p.logger.Error("failed setting embodied metric", "instance", i.Name)
		return
	}

	for _, m := range i.Metrics {
		m := m
		// setup metric labels
		attrs := getAtrributesFromLabels(&m)
		attrs = append(
			attrs,
			attribute.Key("type").String(m.ResourceType.String()),
			attribute.Key("provider").String(i.Provider.String()),
			attribute.Key("type").String(m.ResourceType.String()),
		)

		// register emission metrics for instance
		_, err := p.meter.RegisterCallback(
			func(ctx context.Context, o api.Observer) error {
				o.ObserveFloat64(
					emissions,
					m.Emissions.Value,
					api.WithAttributes(attrs...),
				)
				return nil
			}, emissions)
		if err != nil {
			p.logger.Error("failed setting metric", "instance", i.Name)
		}
	}
}

func getAtrributesFromLabels(m *v1.Metric) []attribute.KeyValue {
	attrs := []attribute.KeyValue{}
	for k, l := range m.Labels {
		attrs = append(attrs, attribute.Key(k).String(l))
	}
	return attrs
}

func getAttributesFromInstance(i *v1.Instance) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.Key("kind").String(i.Kind),
		attribute.Key("name").String(i.Name),
		attribute.Key("zone").String(i.Zone),
		attribute.Key("region").String(i.Region),
		attribute.Key("service").String(i.Service),
		attribute.Key("provider").String(i.Provider.String()),
	}
}
