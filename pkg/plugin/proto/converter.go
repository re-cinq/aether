package proto

import (
	"errors"
	"time"

	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

// ConvertToPB is a helper function that converts the Instance type to its
// protobuffer alternative
func ConvertToPB(src *v1.Instance) (*InstanceRequest, error) {
	if src == nil {
		return nil, errors.New("source struct is nil")
	}

	dest := &InstanceRequest{
		Id:       src.ID,
		Provider: string(src.Provider),
		Service:  src.Service,
		Name:     src.Name,
		Region:   src.Region,
		Zone:     src.Zone,
		Kind:     src.Kind,
		Status:   string(src.Status),
		EmbodiedEmissions: &ResourceEmissions{
			Value: src.EmbodiedEmissions.Value,
			Unit:  string(src.EmbodiedEmissions.Unit),
		},
		Labels: map[string]string(src.Labels),
	}

	metrics := make(map[string]*Metric)
	for key, metric := range src.Metrics {
		metrics[key] = &Metric{
			Name:       metric.Name,
			Usage:      metric.Usage,
			UnitAmount: metric.UnitAmount,
			Unit:       string(metric.Unit),
			Energy:     metric.Energy,
			Emissions: &ResourceEmissions{
				Value: metric.Emissions.Value,
				Unit:  string(metric.Emissions.Unit),
			},
			Labels:    map[string]string(metric.Labels),
			UpdatedAt: metric.UpdatedAt.Unix(),
		}
	}
	dest.Metrics = metrics
	return dest, nil
}

// ConvertToInstance converts a protobugger InstanceRequest struct to a
// v1.Instance type
func ConvertToInstance(src *InstanceRequest) (*v1.Instance, error) {
	if src == nil {
		return nil, errors.New("source struct is nil")
	}

	instance := &v1.Instance{
		ID:       src.Id,
		Provider: v1.Provider(src.Provider),
		Service:  src.Service,
		Name:     src.Name,
		Region:   src.Region,
		Zone:     src.Zone,
		Kind:     src.Kind,
		Status:   v1.InstanceStatus(src.Status),
		Labels:   v1.Labels(src.Labels),
	}

	if src.EmbodiedEmissions != nil {
		instance.EmbodiedEmissions = v1.ResourceEmissions{
			Value: src.EmbodiedEmissions.Value,
			Unit:  v1.EmissionUnit(src.EmbodiedEmissions.Unit),
		}
	}

	instanceMetrics := make(v1.Metrics)
	for key, metric := range src.Metrics {
		instanceMetrics[key] = v1.Metric{
			Name:         metric.Name,
			ResourceType: v1.ResourceType(metric.ResourceType),
			Usage:        metric.Usage,
			Energy:       metric.Energy,
			UnitAmount:   metric.UnitAmount,
			Unit:         v1.ResourceUnit(metric.Unit),
			Emissions: v1.ResourceEmissions{
				Value: metric.Emissions.Value,
				Unit:  v1.EmissionUnit(metric.Emissions.Unit),
			},
			UpdatedAt: time.Unix(metric.UpdatedAt, 0),
			Labels:    v1.Labels(metric.Labels),
		}
	}
	instance.Metrics = instanceMetrics

	return instance, nil
}
