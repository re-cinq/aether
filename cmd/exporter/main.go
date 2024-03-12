package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"log/slog"

	"github.com/re-cinq/cloud-carbon/pkg/api"
	"github.com/re-cinq/cloud-carbon/pkg/calculator"
	"github.com/re-cinq/cloud-carbon/pkg/config"
	"github.com/re-cinq/cloud-carbon/pkg/exporter"
	"github.com/re-cinq/cloud-carbon/pkg/pathfinder"
	"github.com/re-cinq/cloud-carbon/pkg/scheduler"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	bus "github.com/re-cinq/go-bus"
)

const startUpLog = `
  ___  __    _____  __  __  ____      ___    __    ____  ____  _____  _  _ 
 / __)(  )  (  _  )(  )(  )(  _ \    / __)  /__\  (  _ \(  _ \(  _  )( \( )
( (__  )(__  )(_)(  )(__)(  )(_) )  ( (__  /(__)\  )   / ) _ < )(_)(  )  ( 
 \___)(____)(_____)(______)(____/    \___)(__)(__)(_)\_)(____/(_____)(_)\_)
                                                                                              
`

var (
	description    = "Cloud Carbon collection exporter"
	gitSHA         = "n/a"
	name           = "Cloud Carbon"
	source         = "https://github.com/re-cinq/cloud-carbon"
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
	fmt.Printf("Source:         %s\n", source)
	fmt.Printf("Built:          %s\n", buildTimestamp)
}

func main() {
	// Record when the program is started
	start := time.Now()

	ctx := context.Background()

	// print the logo
	slog.Info(startUpLog)

	// check if we got args passed
	args := os.Args

	// Print the version and exit
	if len(args) > 1 && args[1] == "version" {
		PrintVersion()
		return
	}

	// Parse the flags
	flag.Parse()

	// At this point load the config
	config.InitConfig()

	// Init the application bus
	eventBus := bus.NewEventBus(8192, runtime.NumCPU())

	// Subscribe to the metrics collections
	eventBus.Subscribe(
		v1.MetricsCollectedTopic,
		calculator.NewEmissionCalculator(eventBus),
	)

	// Subscribe to update the prometheus exporter
	eventBus.Subscribe(
		v1.EmissionsCalculatedTopic,
		exporter.NewPrometheusEventHandler(eventBus),
	)

	// Subscribe to update the pathfinder handler
	eventBus.Subscribe(
		v1.EmissionsCalculatedTopic,
		pathfinder.NewPathfinderEventHandler(eventBus),
	)

	// Start the bus
	eventBus.Start()

	// Create the API object
	apiServer := api.NewAPIServer()

	// Scheduler manager
	scraper := scheduler.NewScrapingManager(ctx, eventBus)

	// Start the API
	go apiServer.Start()

	// Print the start
	slog.Info("started", "time", time.Since(start))

	// Start the scheduler manager
	scraper.Start(ctx)

	// Graceful shutdown
	// Await for the signals to teminate the program
	await(func() {
		// Shutdown the API server
		apiServer.Stop()

		// Stop all the scraping
		scraper.Stop()

		// Shutdown the bus
		eventBus.Stop()
	})
}

// await for the signals and run the shutdown function
func await(shutdownHook func()) {
	terminating := make(chan bool, 1)

	// Signals channel
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		terminationSignal := <-signalChan

		// Warn that we are terminating
		slog.Info("terminating", "signal", terminationSignal)

		// Run the shutdown hook
		shutdownHook()

		slog.Info("terminated Successfully")

		terminating <- true
	}()

	// Here we wait
	<-terminating
}
