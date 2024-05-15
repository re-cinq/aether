/*
* This is a skeleton plugin that you can use to build an Aether SourcePlugin
* it is minimal with no functionality but adheres to the Aether plugin
* interface
 */
package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"

	aetherplugin "github.com/re-cinq/aether/pkg/plugin"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

// ExampleSource is the struct that will cover the Plugin interface
type ExampleSource struct {
	logger *slog.Logger
}

// Fetch is what gets run by aether, it should return a list of instances with
// the metrics attached to them (CPU, Memory, Networking, Storage)
// Your business logic should be handeled here
func (e *ExampleSource) Fetch(ctx context.Context) (map[string]*v1.Instance, error) {
	// business logic here
	return map[string]*v1.Instance{
		"example": {
			ID:       "123456789",
			Provider: v1.Custom,
			Service:  "On-Prem",
			Status:   v1.InstanceRunning,
			Name:     "vmware1",
			Region:   "europe-west1",
			Zone:     "europe-west1-b",
			Kind:     "e2-medium",
			Metrics: v1.Metrics{
				"cpu": v1.Metric{
					Name:         "cpu",
					ResourceType: "",
					Usage:        0.5949916219267631,
					UnitAmount:   1,
					Unit:         "vCPU",
					Energy:       0.00004006771376096877,
					Emissions: v1.ResourceEmissions{
						Value: 0.005421161671859076,
						Unit:  "gCO2eq",
					},
					UpdatedAt: time.Now(),
					Labels: v1.Labels{
						"id":           "4750323081948662787",
						"machine_type": "e2-medium",
						"name":         "vmare1",
						"region":       "europe-west1",
						"zone":         "europe-west1-b",
					},
				},
			},
			EmbodiedEmissions: v1.ResourceEmissions{
				Value: 0.00005724905663368849,
				Unit:  "gCO2eq",
			},
			Labels: v1.Labels{
				"ID":        "4750323081948662787",
				"Lifecycle": "STANDARD",
			},
		},
	}, nil
}

// Stop is used to adhere to the interface and will be called when Aether is
// shut down
func (e *ExampleSource) Stop(ctx context.Context) error {
	// business logic here
	return nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	// initialize the struct
	example := &ExampleSource{logger}

	// this is important, you will need to set the key to the name of the binary
	// that is built. in this case we use example as we use the following command
	// to build this file
	//
	// ```
	//	GOOS=linux GOARCH=amd64 go build -o .plugins/source/example \
	//		pkg/plugin/example/source/example.go
	// ```
	// A binary is created in the .plugin folder called `example`
	// this directory will be read by aether and all binarys will be started
	// and registered as plugins
	pluginMap := map[string]plugin.Plugin{
		"example": &aetherplugin.SourcePlugin{Impl: example},
	}

	// This is a blocking call that will start the plugin
	plugin.Serve(&plugin.ServeConfig{
		// This should corrspond to the handshake that the version of
		// aether you are using is using
		HandshakeConfig: aetherplugin.SourceHandshake,
		Plugins:         pluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
