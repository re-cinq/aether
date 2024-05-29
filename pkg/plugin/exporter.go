package plugin

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/re-cinq/aether/pkg/log"
	"github.com/re-cinq/aether/pkg/plugin/proto"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

// Handshake is a common handshake that is shared by plugin and host.
// this is to make sure that versioning of plugins is equal
// and is not for security
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "AETHER_EXPORTER_PLUGIN",
	MagicCookieValue: "886c2a46-18b4-4090-8e13-0461439bb0d0",
}

// Exporter is the interface that plugins need to adhere to
type Exporter interface {
	Send(i *v1.Instance) error
}

// GRPCClient is an implementation of Exporter that can communicate over RPC
type GRPCClient struct{ client proto.ExporterClient }

// Send is used to fulfill the Exporter interface
func (g *GRPCClient) Send(i *v1.Instance) error {
	ir, err := proto.ConvertToPB(i)

	if err != nil {
		return fmt.Errorf("failed converting to protobuf, %v", err)
	}

	_, err = g.client.Send(context.Background(), ir)
	return err
}

// GRPCServer is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	Impl Exporter
}

// Send is used to receive the RPC requests
func (m *GRPCServer) Send(
	ctx context.Context,
	req *proto.InstanceRequest,
) (*proto.Empty, error) {
	i, err := proto.ConvertToInstance(req)
	if err != nil {
		return &proto.Empty{}, err
	}
	return &proto.Empty{}, m.Impl.Send(i)
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type ExporterPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl Exporter
}

// GRPCServer is used to setup a GRPC server
func (p *ExporterPlugin) GRPCServer(
	broker *plugin.GRPCBroker,
	s *grpc.Server,
) error {
	proto.RegisterExporterServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

// GRPCClient returns a gRPC client
func (p *ExporterPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &GRPCClient{client: proto.NewExporterClient(c)}, nil
}

// ExportPluginSystem is used to manage all third party exporters
type ExportPluginSystem struct {
	Dir string

	Plugins []*RegisteredPlugin
}

// RegisteredPlugin is all the information needed to manage a registered plugin
type RegisteredPlugin struct {
	Name     string
	Exporter Exporter
	Client   plugin.ClientProtocol
}

// Load goes to a directory and runs every binary in that directory as a plugin
// The plugins for the exporter will need to adhere to the ExporterPlugin
// interface else they will fail to load
func (e *ExportPluginSystem) Load(ctx context.Context) error {
	logger := log.FromContext(ctx)

	// Open the plugin directory.
	files, err := os.ReadDir(e.Dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}
		// Run the plugin executable as a subprocess.
		pluginPath := filepath.Join(e.Dir, file.Name())

		// Connect to the plugin process over RPC.
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: Handshake,
			Plugins: map[string]plugin.Plugin{
				file.Name(): &ExporterPlugin{},
			},
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
			Cmd:              exec.Command(pluginPath),
			Logger: hclog.New(&hclog.LoggerOptions{
				Name:       "exporter",
				Output:     os.Stdout,
				Level:      hclog.Debug,
				JSONFormat: true,
			}),
		})

		// Start the plugin process.
		rpcClient, err := client.Client()
		if err != nil {
			logger.Error("connecting to plugin process", "error", err)
			continue
		}

		// Request the plugin.
		raw, err := rpcClient.Dispense(file.Name())
		if err != nil {
			logger.Error("dispensing plugin", "error", err)
			continue
		}

		p := raw.(Exporter)

		// make sure plugin is reachable
		err = rpcClient.Ping()
		if err != nil {
			logger.Info("failed loading plugin", "plugin", file.Name(), "error", err)
		}

		e.Plugins = append(e.Plugins, &RegisteredPlugin{
			Name:     file.Name(),
			Exporter: p,
			Client:   rpcClient,
		})
		logger.Info("loaded plugin", "plugin", file.Name())
	}

	return nil
}
