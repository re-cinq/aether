package proto

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

func TestConvertToPB(t *testing.T) {
	tests := []struct {
		name     string
		src      *v1.Instance
		expected *InstanceRequest
	}{
		{
			name: "Valid instance",
			src: &v1.Instance{
				ID:                "1",
				Provider:          "test",
				Service:           "test-service",
				Name:              "test-instance",
				Region:            "test-region",
				Zone:              "test-zone",
				Kind:              "test-kind",
				EmbodiedEmissions: v1.ResourceEmissions{Value: 200, Unit: "test-unit"},
				Labels:            map[string]string{"label1": "value1", "label2": "value2"},
				Metrics: map[string]v1.Metric{
					"metric1": {
						Name:       "metric1",
						Usage:      100,
						UnitAmount: 10.5,
						Unit:       "test-unit",
						Emissions:  v1.ResourceEmissions{Value: 50, Unit: "test-unit"},
						Labels:     map[string]string{"label1": "value1"},
						UpdatedAt:  time.Now(),
					},
				},
			},
			expected: &InstanceRequest{
				Id:                "1",
				Provider:          "test",
				Service:           "test-service",
				Name:              "test-instance",
				Region:            "test-region",
				Zone:              "test-zone",
				Kind:              "test-kind",
				EmbodiedEmissions: &ResourceEmissions{Value: 200, Unit: "test-unit"},
				Labels:            map[string]string{"label1": "value1", "label2": "value2"},
				Metrics: map[string]*Metric{
					"metric1": {
						Name:       "metric1",
						Usage:      100,
						UnitAmount: 10.5,
						Unit:       "test-unit",
						Emissions:  &ResourceEmissions{Value: 50, Unit: "test-unit"},
						Labels:     map[string]string{"label1": "value1"},
						UpdatedAt:  time.Now().Unix(),
					},
				},
			},
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := ConvertToPB(test.src)
			if err != nil {
				t.Fatalf("ConvertToPB failed with error: %v", err)
			}

			// Compare result with expected
			if !compareInstanceRequests(result, test.expected) {
				t.Errorf("ConvertToPB result does not match expected")
			}
		})
	}
}

func compareInstanceRequests(a, b *InstanceRequest) bool {
	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}
	if a == nil && b == nil {
		return true
	}

	// Compare exported fields
	if a.Provider != b.Provider || a.Service != b.Service || a.Name != b.Name ||
		a.Region != b.Region || a.Zone != b.Zone || a.Kind != b.Kind || a.Id != b.Id {
		return false
	}

	// Compare EmbodiedEmissions
	if !compareResourceEmissions(a.EmbodiedEmissions, b.EmbodiedEmissions) {
		return false
	}

	// Compare Labels
	if !compareStringMaps(a.Labels, b.Labels) {
		return false
	}

	// Compare Metrics
	if !compareMetricsMap(a.Metrics, b.Metrics) {
		return false
	}

	return true
}

func compareResourceEmissions(a, b *ResourceEmissions) bool {
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}
	if a == nil && b == nil {
		return true
	}

	return a.Value == b.Value && a.Unit == b.Unit
}

func compareStringMaps(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	for key, valueA := range a {
		valueB, ok := b[key]
		if !ok || valueA != valueB {
			return false
		}
	}

	return true
}

func compareMetricsMap(a, b map[string]*Metric) bool {
	if len(a) != len(b) {
		return false
	}

	for key, metricA := range a {
		metricB, ok := b[key]
		if !ok || !compareMetrics(metricA, metricB) {
			return false
		}
	}

	return true
}

func compareMetrics(a, b *Metric) bool {
	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}
	if a == nil && b == nil {
		return true
	}

	return a.Name == b.Name && a.Usage == b.Usage &&
		a.UnitAmount == b.UnitAmount && a.Unit == b.Unit &&
		a.UpdatedAt == b.UpdatedAt && compareStringMaps(a.Labels, b.Labels) &&
		compareResourceEmissions(a.Emissions, b.Emissions)
}

func TestConvertToInstance(t *testing.T) {
	tests := []struct {
		name     string
		src      *InstanceRequest
		expected *v1.Instance
		err      error
	}{
		{
			name: "Valid InstanceRequest",
			src: &InstanceRequest{
				Id:       "1",
				Provider: "test",
				Service:  "test-service",
				Name:     "test-instance",
				Region:   "test-region",
				Zone:     "test-zone",
				Kind:     "test-kind",
				EmbodiedEmissions: &ResourceEmissions{
					Value: 200,
					Unit:  "test-unit",
				},
				Labels: map[string]string{"label1": "value1", "label2": "value2"},
				Metrics: map[string]*Metric{
					"metric1": {
						Name:         "metric1",
						ResourceType: "test-resource-type",
						Usage:        100,
						UnitAmount:   10.5,
						Unit:         "test-unit",
						Emissions: &ResourceEmissions{
							Value: 50,
							Unit:  "test-unit",
						},
						UpdatedAt: 1234567890,
						Labels:    map[string]string{"label1": "value1"},
					},
				},
			},
			expected: &v1.Instance{
				ID:       "1",
				Provider: "test",
				Service:  "test-service",
				Name:     "test-instance",
				Region:   "test-region",
				Zone:     "test-zone",
				Kind:     "test-kind",
				Labels:   v1.Labels{"label1": "value1", "label2": "value2"},
				EmbodiedEmissions: v1.ResourceEmissions{
					Value: 200,
					Unit:  v1.EmissionUnit("test-unit"),
				},
				Metrics: v1.Metrics{
					"metric1": {
						Name:         "metric1",
						ResourceType: v1.ResourceType("test-resource-type"),
						Usage:        100,
						UnitAmount:   10.5,
						Unit:         v1.ResourceUnit("test-unit"),
						Emissions: v1.ResourceEmissions{
							Value: 50,
							Unit:  v1.EmissionUnit("test-unit"),
						},
						UpdatedAt: time.Unix(1234567890, 0),
						Labels:    v1.Labels{"label1": "value1"},
					},
				},
			},
			err: nil,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := ConvertToInstance(test.src)

			if !errors.Is(err, test.err) {
				t.Errorf("Expected error %v, got %v", test.err, err)
			}

			if diff := cmp.Diff(test.expected, result); diff != "" {
				t.Errorf("Mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
