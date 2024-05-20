package calculator

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"

	"gopkg.in/yaml.v2"

	"github.com/re-cinq/aether/pkg/bus"
	"github.com/re-cinq/aether/pkg/config"
	"github.com/re-cinq/aether/pkg/log"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
	factors "github.com/re-cinq/aether/pkg/types/v1/factors"
	data "github.com/re-cinq/emissions-data/pkg/types/v2"
)

var instanceData map[string]data.Instance

// AWS, GCP and Azure have increased their server lifespan to 6 years (2024)
// https://sustainability.aboutamazon.com/products-services/the-cloud?energyType=true
// https://www.theregister.com/2024/01/31/alphabet_q4_2023/
// https://www.theregister.com/2022/08/02/microsoft_server_life_extension/
const serverLifespan = 6

//nolint:govet // don't need to write the struct fields each time
var emptyWattage = []data.Wattage{{0, 0}, {10, 0}, {50, 0}, {100, 0}}

// CalculatorHandler is used to handle events when metrics have been collected
type CalculatorHandler struct {
	Bus    *bus.Bus
	logger *slog.Logger
}

// NewHandler returns a new configuered instance of CalculatorHandler
// as well as setups the factor datasets
func NewHandler(ctx context.Context, b *bus.Bus) *CalculatorHandler {
	logger := log.FromContext(ctx)

	err := factors.CloneAndUpdateFactorsData()
	if err != nil {
		logger.Error("error with emissions repo", "error", err)
		return nil
	}

	for provider := range config.AppConfig().Providers {
		err = getProviderEmissionFactors(provider)
		if err != nil {
			logger.Error("unable to get v2 Emission Factors", "error", err, "provider", provider)
		}
	}

	return &CalculatorHandler{
		Bus:    b,
		logger: logger,
	}
}

// Stop is used to fulfill the EventHandler interface and all clean up
// functionality should be run in here
func (c *CalculatorHandler) Stop(ctx context.Context) {}

// Handle is used to fulfill the EventHandler interface and recives an event
// when handler is subscribed to it. Currently only handles v1.MetricsCollectedEvent
func (c *CalculatorHandler) Handle(ctx context.Context, e *bus.Event) {
	switch e.Type {
	case v1.MetricsCollectedEvent:
		c.handleEvent(e)
	default:
		return
	}
}

// handleEvent is the business logic for handeling a v1.MetricsCollectedEvent
// and runs the emissions calculations on the metrics that where received
func (c *CalculatorHandler) handleEvent(e *bus.Event) {
	interval := config.AppConfig().ProvidersConfig.Interval
	ctx := log.WithContext(context.Background(), c.logger)

	instance, ok := e.Data.(v1.Instance)
	if !ok {
		c.logger.Error("EmissionCalculator got an unknown event", "event", e)
		return
	}

	// if an instance is terminated we do not need to calculate
	// emissions for it
	if instance.Status == v1.InstanceTerminated {
		if err := c.Bus.Publish(&bus.Event{
			Type: v1.EmissionsCalculatedEvent,
			Data: instance,
		}); err != nil {
			c.logger.Error("failed publishing terminated instance", "instance", instance.Name, "error", err)
		}
		return
	}

	// Gets PUE, grid data, and machine specs
	factor, err := factors.ProviderEmissions(instance.Provider, factors.DataPath)
	if err != nil {
		c.logger.Error("error getting emission factors", "error", err)
		return
	}

	grid, ok := factor.Coefficient[instance.Region]
	if !ok {
		c.logger.Error("region not found in factors", "region", instance.Region, "provider", instance.Provider)
		return
	}
	// TODO: hotfix until updated in emissions data
	// convert gridCO2e from metric tonnes to grams
	grid *= (1000 * 1000)

	// energy consumption is already computed for kepler, so we just need to calculate
	// the carbon emissions
	// TODO this is a hack until the energy consumption is calculator is separated from
	// the emissions calculator
	if instance.Service == "kepler" {
		for _, m := range instance.Metrics {
			m.Emissions = v1.NewResourceEmission(
				m.Energy*factor.AveragePUE*grid,
				v1.GCO2eq,
			)
			if err != nil {
				c.logger.Error("error calulating emissions", "type", m.Name, "error", err)
				continue
			}
			// update the instance metrics
			instance.Metrics.Upsert(&m)
		}

		// We publish the interface on the bus once its been calculated
		if err := c.Bus.Publish(&bus.Event{
			Type: v1.EmissionsCalculatedEvent,
			Data: instance,
		}); err != nil {
			c.logger.Error("failed publishing instance after calculation", "instance", instance.Name, "error", err)
		}
		return
	}

	params := &parameters{
		grid: grid,
		pue:  factor.AveragePUE,
	}

	specs, ok := factor.Embodied[instance.Kind]
	if !ok {
		c.logger.Error("failed finding instance in factor data", "instance", instance.Name, "kind", instance.Kind)
		return
	}

	if d, ok := instanceData[instance.Kind]; ok {
		params.factors = &d
	}

	// fallback to use spec power min and max watt values.
	// this is less accurate and a place holder until a
	// different solution is implemented.
	if reflect.DeepEqual(params.factors.PkgWatt, emptyWattage) {
		params.factors.PkgWatt = []data.Wattage{
			{
				Percentage: 0,
				Wattage:    specs.MinWatts,
			},
			{
				Percentage: 100,
				Wattage:    specs.MaxWatts,
			},
		}
	}

	if params.embodiedFactor == 0 {
		params.embodiedFactor = hourlyEmbodiedEmissions(&specs)
	}

	// calculate and set the operational emissions for each
	// metric type (CPU, Memory, Storage, and networking)
	metrics := instance.Metrics
	for _, v := range metrics {
		params.metric = &v

		err := operationalEmissions(ctx, interval, params)
		if err != nil {
			c.logger.Error("error calulating emissions", "type", v.Name, "error", err)
			continue
		}
		// update the instance metrics
		metrics.Upsert(params.metric)
	}

	instance.EmbodiedEmissions = v1.NewResourceEmission(
		embodiedEmissions(interval, params.embodiedFactor),
		v1.GCO2eq,
	)

	// We publish the interface on the bus once its been calculated
	if err := c.Bus.Publish(&bus.Event{
		Type: v1.EmissionsCalculatedEvent,
		Data: instance,
	}); err != nil {
		c.logger.Error("failed publishing instance after calculation", "instance", instance.Name, "error", err)
	}
}

func hourlyEmbodiedEmissions(e *factors.Embodied) float64 {
	// we fall back on the specs from the previous dataset
	// and convert it into a hourly factor
	// this is based on CCF's calculation:
	//
	// M = TE * (TR/EL) * (RR/TR)
	//
	// TE = Total Embodied Emissions
	// TR = Time Reserved (in years)
	// EL = Expected Lifespan
	// RR = Resources Reserved
	// TR = Total Resources, the total number of resources available.
	return e.TotalEmbodiedKiloWattCO2e *
		// 1 hour normalized to a year
		((1.0 / 24.0 / 365.0) / serverLifespan) *
		// amount of vCPUS for instance versus total vCPUS for platform
		(e.VCPU / e.TotalVCPU)
}

func getProviderEmissionFactors(provider v1.Provider) error {
	url := "https://raw.githubusercontent.com/re-cinq/emissions-data/main/data/v2/%s-instances.yaml"
	u := fmt.Sprintf(url, provider)

	r, err := http.Get(u)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	d := yaml.NewDecoder(r.Body)
	err = d.Decode(&instanceData)
	if err != nil {
		return err
	}

	return nil
}
