package exporter

import (
	"context"

	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	bus "github.com/re-cinq/go-bus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"k8s.io/klog/v2"
)

type PrometheusEventHandler struct {
	eventBus bus.Bus
	meter    api.Meter
}

func NewPrometheusEventHandler(eventBus bus.Bus) *PrometheusEventHandler {
	exporter, err := prometheus.New()
	if err != nil {
		klog.Fatal(err)
	}

	meter := metric.NewMeterProvider(
		metric.WithReader(exporter),
	).Meter("cloud-carbon")

	return &PrometheusEventHandler{
		eventBus: eventBus,
		meter:    meter,
	}
}

func (c *PrometheusEventHandler) Apply(event bus.Event) {
	e, ok := event.(v1.EmissionsCalculated)
	if !ok {
		klog.Errorf("PrometheusEventHandler got an unknown event: %T", event)
		return
	}

	// setup emissions gauge
	emissions, err := c.meter.Float64ObservableGauge(
		"emissions",
		api.WithDescription("co2eq of various services"),
	)
	if err != nil {
		klog.Errorf("[otel] failed setting up emissions metric")
		return
	}

	// setup embodied emissions gauge
	embodied, err := c.meter.Float64ObservableGauge(
		"embodied",
		api.WithDescription("co2eq of various services"),
	)
	if err != nil {
		klog.Errorf("[otel] failed setting up embodied emissions metric")
		return
	}

	// register embodied emissions metrics for instance
	// NOTE: this will not change based on different types of metrics
	_, err = c.meter.RegisterCallback(
		func(ctx context.Context, o api.Observer) error {
			o.ObserveFloat64(
				embodied,
				e.Instance.EmbodiedEmissions().Value(),
				api.WithAttributes(
					getAttributesFromInstance(&e.Instance)...,
				))

			return nil
		}, embodied)
	if err != nil {
		klog.Errorf("error occurred setting embodied metric for %v", e.Instance.Name())
		return
	}

	for _, m := range e.Instance.Metrics() {
		m := m
		// setup metric labels
		attrs := getAtrributesFromLabels(&m)
		attrs = append(
			attrs,
			attribute.Key("type").String(m.Type().String()),
			attribute.Key("provider").String(e.Instance.Provider().String()),
		)

		// register emission metrics for instance
		_, err := c.meter.RegisterCallback(
			func(ctx context.Context, o api.Observer) error {
				o.ObserveFloat64(
					emissions,
					m.Emissions().Value(),
					api.WithAttributes(attrs...),
				)
				return nil
			}, emissions)
		if err != nil {
			klog.Errorf("error occurred setting metric for %v", e.Instance.Name())
		}
	}
}

func getAtrributesFromLabels(m *v1.Metric) []attribute.KeyValue {
	attrs := []attribute.KeyValue{}
	for k, l := range m.Labels() {
		attrs = append(attrs, attribute.Key(k).String(l))
	}
	return attrs
}

func getAttributesFromInstance(i *v1.Instance) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.Key("kind").String(i.Kind()),
		attribute.Key("name").String(i.Name()),
		attribute.Key("zone").String(i.Zone()),
		attribute.Key("region").String(i.Region()),
		attribute.Key("service").String(i.Service()),
		attribute.Key("provider").String(i.Provider().String()),
	}
}
