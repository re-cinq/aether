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

// SourceHandshake is a common handshake that is shared by plugin and host.
// this is to make sure that versioning of plugins is equal
// and is not for security
var SourceHandshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "AETHER_SOURCE_PLUGIN",
	MagicCookieValue: "9cf50efe-f360-4c46-997f-e1ce7317adaf",
}

// SourceGRPCClient is an implemntation of v1.Source that can
// communicate over RPC
type SourceGRPCClient struct{ client proto.SourceClient }

// Fetch is used to fullfil the v1.Source interface over RPC
func (g *SourceGRPCClient) Fetch(ctx context.Context) ([]*v1.Instance, error) {
	var instances []*v1.Instance
	res, err := g.client.Fetch(ctx, &proto.Empty{})

	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, fmt.Errorf("no instances returned from plugin")
	}

	for i := range res.Instances {
		r, convertErr := proto.ConvertToInstance(res.Instances[i])
		if err != nil {
			return nil, convertErr
		}
		instances = append(instances, r)
	}

	return instances, err
}

// Stop is used to fullfil the v1.Source interface over RPC
func (g *SourceGRPCClient) Stop(ctx context.Context) error {
	return nil
}

// SourceGRPCServer is an GRPC implementation over which
// SourceGRPClient communicates
type SourceGRPCServer struct {
	Impl v1.Source
}

// Fetch is used to fulfill the GRPC Interface
func (s *SourceGRPCServer) Fetch(
	ctx context.Context,
	p *proto.Empty,
) (*proto.ListInstanceResponse, error) {
	var listRes proto.ListInstanceResponse
	res, err := s.Impl.Fetch(ctx)

	// at this point im not sure how to better deal with this
	// as we need to convert each Protobuf struct to a v1.Instance struct
	for i := range res {
		r, _ := proto.ConvertToPB(res[i])
		listRes.Instances = append(listRes.Instances, r)
	}
	return &listRes, err
}

// Stop is used to fulfill the GRPC Interface
func (s *SourceGRPCServer) Stop(
	ctx context.Context,
	p *proto.Empty,
) (*proto.Empty, error) {
	return &proto.Empty{}, s.Impl.Stop(ctx)
}

// SourcePlugin is the implementation of plugin.GRPCPlugin so we can cionsume
// this plugin
type SourcePlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl v1.Source
}

// GRPCServer is used to setup a GRPC server
func (s *SourcePlugin) GRPCServer(
	broker *plugin.GRPCBroker,
	srv *grpc.Server,
) error {
	proto.RegisterSourceServer(srv, &SourceGRPCServer{Impl: s.Impl})
	return nil
}

// GRPCClient returns a gRPC client
func (s *SourcePlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &SourceGRPCClient{client: proto.NewSourceClient(c)}, nil
}

// SourcePluginSystem is used to manage all third party source
type SourcePluginSystem struct {
	Dir string

	Plugins []*RegisteredSourcePlugin
}

// RegisteredSourcePlugin is all the information needed to manage a registered plugin
type RegisteredSourcePlugin struct {
	Name   string
	Source v1.Source
	Client plugin.ClientProtocol
}

// Load goes to a directory and runs every binary in that directory as a plugin
// The plugins for the sources will need to adhere to the SourcePlugin
// interface else they will fail to load
func (s *SourcePluginSystem) Load(ctx context.Context) error {
	logger := log.FromContext(ctx)

	// Open the plugin directory.
	files, err := os.ReadDir(s.Dir)
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
		pluginPath := filepath.Join(s.Dir, file.Name())

		// Connect to the plugin process over RPC.
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: SourceHandshake,
			Plugins: map[string]plugin.Plugin{
				file.Name(): &SourcePlugin{},
			},
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
			Cmd:              exec.Command(pluginPath),
			Logger: hclog.New(&hclog.LoggerOptions{
				Name:       "source",
				Output:     os.Stdout,
				Level:      getLogLevel(),
				JSONFormat: true,
			}),
		})

		// Start the plugin process.
		rpcClient, err := client.Client()
		if err != nil {
			logger.Error("connecting to source plugin process", "error", err)
			continue
		}

		// Request the plugin.
		raw, err := rpcClient.Dispense(file.Name())
		if err != nil {
			logger.Error("dispensing plugin", "error", err)
			continue
		}

		p := raw.(v1.Source)

		// make sure plugin is reachable
		err = rpcClient.Ping()
		if err != nil {
			logger.Info("failed loading plugin", "plugin", file.Name(), "error", err)
		}

		s.Plugins = append(s.Plugins, &RegisteredSourcePlugin{
			Name:   file.Name(),
			Source: p,
			Client: rpcClient,
		})
		logger.Info("loaded source plugin", "plugin", file.Name())
	}

	return nil
}
