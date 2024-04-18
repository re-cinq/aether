/*
* This is a skeleton plugin that you can use to build an Aether ExporterPlugin
* it is minimal with no functionality but adheres to the Aether plugin
* interface
 */
package main

import (
	"log/slog"
	"os"

	"github.com/hashicorp/go-plugin"

	aetherplugin "github.com/re-cinq/aether/pkg/plugin"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

// ExampleExporteris the struct that will cover the Plugin interface
type ExampleExporter struct {
	logger *slog.Logger
}

// Send is what gets run by aether, it should receive a
// v1.Instance that has metrics with the emissions already calculated
// Your business logic should be handeled here
func (e *ExampleExporter) Send(i *v1.Instance) error {
	e.logger.Info("received instance", "instance", &i)

	// business logic here
	return nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	// initialize the struct
	example := &ExampleExporter{logger}

	// this is important, you will need to set the key to the name of the binary
	// that is built. in this case we use example as we use the following command
	// to build this file
	//
	// ```
	//	GOOS=linux GOARCH=amd64 go build -o .plugins/example \
	//		pkg/plugin/example/example.go
	// ```
	// A binary is created in the .plugin folder called `example`
	// this directory will be read by aether and all binarys will be started
	// and registered as plugins
	pluginMap := map[string]plugin.Plugin{
		"example": &aetherplugin.ExporterPlugin{Impl: example},
	}

	// This is a blocking call that will start the plugin
	plugin.Serve(&plugin.ServeConfig{
		// This should corrspond to the handshake that the version of
		// aether you are using is using
		HandshakeConfig: aetherplugin.Handshake,
		Plugins:         pluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
