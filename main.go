package main

import (
	"context"
	"flag"
	"os"
	"runtime"
	"time"

	"github.com/re-cinq/cloud-carbon/pkg/api"
	"github.com/re-cinq/cloud-carbon/pkg/calculator"
	"github.com/re-cinq/cloud-carbon/pkg/config"
	"github.com/re-cinq/cloud-carbon/pkg/exporter"
	"github.com/re-cinq/cloud-carbon/pkg/pathfinder"
	"github.com/re-cinq/cloud-carbon/pkg/scheduler"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	bus "github.com/re-cinq/go-bus"
	"k8s.io/klog/v2"
)

const startUpLog = `
  ___  __    _____  __  __  ____      ___    __    ____  ____  _____  _  _ 
 / __)(  )  (  _  )(  )(  )(  _ \    / __)  /__\  (  _ \(  _ \(  _  )( \( )
( (__  )(__  )(_)(  )(__)(  )(_) )  ( (__  /(__)\  )   / ) _ < )(_)(  )  ( 
 \___)(____)(_____)(______)(____/    \___)(__)(__)(_)\_)(____/(_____)(_)\_)
                                                                                              
`

func main() {
	// Record when the program is started
	start := time.Now()

	ctx := context.Background()

	// print the logo
	klog.Infof("%v", startUpLog)

	// check if we got args passed
	args := os.Args

	// Print the version and exit
	if len(args) > 1 && args[1] == "version" {
		PrintVersion()
		return
	}

	// Parse the flags
	klog.InitFlags(nil)
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
	klog.Infof("started in %v", time.Since(start))

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
