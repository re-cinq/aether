package main

import (
	"os"
	"os/signal"
	"syscall"

	"k8s.io/klog/v2"
)

// await for the signals and run the shutdown function
func await(shutdownHook func()) {

	terminating := make(chan bool, 1)

	// Signals channel
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		terminationSignal := <-signalChan

		// Warn that we are terminating
		klog.Infof("terminating: %s", terminationSignal)

		// Run the shutdown hook
		shutdownHook()

		klog.Info("--- terminated Successfully")

		terminating <- true
	}()

	// Here we wait
	<-terminating

}
