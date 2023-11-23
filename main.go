package main

import (
	"flag"
	"os"
	"runtime"
	"time"

	"github.com/re-cinq/cloud-carbon/pkg/api"
	"github.com/re-cinq/cloud-carbon/pkg/bus"
	"github.com/re-cinq/cloud-carbon/pkg/calculator"
	"github.com/re-cinq/cloud-carbon/pkg/config"
	"github.com/re-cinq/cloud-carbon/pkg/exporter"
	"github.com/re-cinq/cloud-carbon/pkg/pathfinder"
	"github.com/re-cinq/cloud-carbon/pkg/scheduler"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	"k8s.io/klog/v2"
)

const startUpLog = `
  ___  __    _____  __  __  ____      ___    __    ____  ____  _____  _  _ 
 / __)(  )  (  _  )(  )(  )(  _ \    / __)  /__\  (  _ \(  _ \(  _  )( \( )
( (__  )(__  )(_)(  )(__)(  )(_) )  ( (__  /(__)\  )   / ) _ < )(_)(  )  ( 
 \___)(____)(_____)(______)(____/    \___)(__)(__)(_)\_)(____/(_____)(_)\_)
                                                                                              
`

// var (
// 	metricsPath = flag.String("metrics-path", "/metrics", "metrics path")
// 	port        = flag.String("port", "8000", "")
// 	host        = flag.String("host", "127.0.0.1", "")
// )

func main() {

	// Record when the program is started
	start := time.Now()

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

	// Init the application bus
	eventBus := bus.NewEventBus(8192, runtime.NumCPU())

	// Subscribe to the metrics collections
	eventBus.Subscribe(v1.MetricsCollectedTopic, calculator.NewEmissionCalculator(eventBus))

	// Subscribe to update the prometheus exporter
	eventBus.Subscribe(v1.MetricsCollectedTopic, exporter.NewPrometheusEventHandler(eventBus))

	// Subscribe to update the pathfinder handler
	eventBus.Subscribe(v1.MetricsCollectedTopic, pathfinder.NewPahfinderEventHandler(eventBus))

	// Start the bus
	eventBus.Start()

	// At this point load the config
	config.InitConfig()

	// Scheduler manager

	scheduler := scheduler.NewScrapingManager(eventBus)
	scheduler.Start()

	// Create the API object
	apiServer := api.NewApiServer()

	// Start the API
	go apiServer.Start()

	// Print the start
	klog.Infof("started in %v\n%v", time.Since(start), startUpLog)

	// Graceful shutdown
	// Await for the signals to teminate the program
	await(func() {
		// Shutdown the API server
		apiServer.Stop()

		// Stop all the scraping
		scheduler.Stop()

		// Shutdown the bus
		eventBus.Stop()
	})

}
