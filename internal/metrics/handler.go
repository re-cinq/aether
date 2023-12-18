package metrics

import (
	"context"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"k8s.io/klog/v2"
)

// Handler is a struct that is used to
// setup and handle the metrics exporter
type Handler struct {
	MeterName string
	Meter     api.Meter
	Exporter  prometheus.Exporter
}

// New returns a new instance of the metric handler
func New(name string) *Handler {
	exporter, err := prometheus.New()
	if err != nil {
		klog.Fatal(err)
	}
	return &Handler{
		Exporter: *exporter,
		Meter:    metric.NewMeterProvider(metric.WithReader(exporter)).Meter(name),
	}
}

// Setup intilizizes all metrics
func (h *Handler) Setup(ctx context.Context) error {
	// setup metrics
	gauge, err := h.Meter.Float64ObservableGauge(
		"emissions_cpu",
		api.WithDescription("cpu emissions"),
	)
	if err != nil {
		return err
	}

	// handle metrics
	err = h.ExampleHandler(ctx, gauge)
	return err
}

// ExampleHandler is an example of how we would set the value of a
// gauge metric
func (h *Handler) ExampleHandler(
	ctx context.Context,
	g api.Float64ObservableGauge,
) error {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, v := range []string{"instance-1", "instance-2"} {
		i := v
		// every time the metrics endpoint is hit it will call this function
		_, err := h.Meter.RegisterCallback(
			func(ctx context.Context, o api.Observer) error {
				// observe a random value
				o.ObserveFloat64(g, rng.Float64()*(100), api.WithAttributes(
					// these set labels on the metric
					attribute.Key("service").String(i),
				))

				return nil
			}, g)
		if err != nil {
			return err
		}
	}

	return nil
}
