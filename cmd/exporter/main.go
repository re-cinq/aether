package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"log/slog"

	"github.com/re-cinq/aether/pkg/api"
	"github.com/re-cinq/aether/pkg/bus"
	"github.com/re-cinq/aether/pkg/calculator"
	"github.com/re-cinq/aether/pkg/config"
	"github.com/re-cinq/aether/pkg/exporter"
	"github.com/re-cinq/aether/pkg/log"
	"github.com/re-cinq/aether/pkg/plugin"
	"github.com/re-cinq/aether/pkg/source"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

const shutdownTTL = time.Second * 15

var (
	description    = "Aether collection exporter"
	gitSHA         = "n/a"
	name           = "Aether"
	repo           = "https://github.com/re-cinq/aether"
	version        = "0.0.1-dev"
	refType        = "branch" // branch or tag
	refName        = ""       // the name of the branch or tag
	buildTimestamp = ""
)

func PrintVersion() {
	fmt.Printf("Name:           %s\n", name)
	fmt.Printf("Version:        %s\n", version)
	fmt.Printf("RefType:        %s\n", refType)
	fmt.Printf("RefName:        %s\n", refName)
	fmt.Printf("Git Commit:     %s\n", gitSHA)
	fmt.Printf("Description:    %s\n", description)
	fmt.Printf("Go Version:     %s\n", runtime.Version())
	fmt.Printf("OS / Arch:      %s / %s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Source:         %s\n", repo)
	fmt.Printf("Built:          %s\n", buildTimestamp)
}

func main() {
	logger, lvl := setupLogger()

	// Record when the program is started
	start := time.Now()

	ctx := context.Background()

	// add logger to context
	ctx = log.WithContext(ctx, logger)

	// check if we got args passed
	args := os.Args

	// Print the version and exit
	if len(args) > 1 && args[1] == "version" {
		PrintVersion()
		return
	}

	// At this point load the config
	config.InitConfig(ctx)

	setLogLevel(lvl, config.AppConfig().LogLevel)

	// Load exporter plugins
	logger.Info("loading plugins")
	pluginsystem := &plugin.ExportPluginSystem{
		Dir: config.AppConfig().Plugins.ExporterDir,
	}

	err := pluginsystem.Load(ctx)

	if err != nil {
		logger.Error("failed to load exporter plugin system", "error", err)
		os.Exit(1)
	}

	// Init the application bus
	b := bus.New()

	// Subscribe to the metrics collections
	b.Subscribe(
		v1.MetricsCollectedEvent,
		calculator.NewHandler(ctx, b),
	)

	// Subscribe to update the prometheus exporter
	b.Subscribe(
		v1.EmissionsCalculatedEvent,
		exporter.NewHandler(ctx, b),
	)

	b.Subscribe(
		v1.EmissionsCalculatedEvent,
		plugin.NewHandler(ctx, pluginsystem),
	)

	// Start the bus
	b.Start(ctx)
	logger.Info("bus started")

	// Create the API object
	server := api.New(api.WithExportPluginSystem(pluginsystem))

	// Scheduler manager
	sourceManager := source.New(ctx, b)
	sourceManager.Start(ctx)
	logger.Info("sources loaded")

	// Start the API
	go server.Start(ctx)

	// Print the start
	logger.Info("started", "time", time.Since(start))

	// Graceful shutdown
	// Await for the signals to teminate the program
	await(ctx, func() {
		// Create a timeout context
		// We need to expect that processes will shutdown in this amount of time
		// or we need to force them
		cancelCtx, cancel := context.WithTimeout(ctx, shutdownTTL)
		defer cancel()

		// Shutdown the API server
		server.Stop(cancelCtx)

		// deregister sources
		sourceManager.Stop(ctx)

		// Shutdown the bus
		b.Stop(ctx)
	})
}

// await for the signals and run the shutdown function
func await(ctx context.Context, shutdown func()) {
	logger := log.FromContext(ctx)
	terminated := make(chan bool, 1)

	// Signals channel
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		terminationSignal := <-signalChan

		// Warn that we are terminating
		logger.Info("terminating", "signal", terminationSignal)

		// Run the shutdown hook
		shutdown()

		terminated <- true
	}()

	// Here we wait
	<-terminated
	logger.Info("terminated successfully")
}

func setupLogger() (*slog.Logger, *slog.LevelVar) {
	// the default log level is INFO
	lvl := new(slog.LevelVar)
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})), lvl
}

func setLogLevel(lvl *slog.LevelVar, level string) {
	switch level {
	case "debug":
		lvl.Set(slog.LevelDebug)
	case "info":
		lvl.Set(slog.LevelInfo)
	case "warn":
		lvl.Set(slog.LevelWarn)
	case "error":
		lvl.Set(slog.LevelError)
	default:
		lvl.Set(slog.LevelInfo)
	}
}
