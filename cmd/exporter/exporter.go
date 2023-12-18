package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/re-cinq/cloud-carbon/internal/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog/v2"
)

func healthProbe(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`ok`))
	if err != nil {
		klog.Fatalf("%s", fmt.Sprintf("failed to write response: %v", err))
	}
}

const startUpLog = `
  ___  __    _____  __  __  ____      ___    __    ____  ____  _____  _  _ 
 / __)(  )  (  _  )(  )(  )(  _ \    / __)  /__\  (  _ \(  _ \(  _  )( \( )
( (__  )(__  )(_)(  )(__)(  )(_) )  ( (__  /(__)\  )   / ) _ < )(_)(  )  ( 
 \___)(____)(_____)(______)(____/    \___)(__)(__)(_)\_)(____/(_____)(_)\_)
                                                                                              
`

var (
	metricsPath = flag.String("metrics-path", "/metrics", "metrics path")
	port        = flag.String("port", "8000", "")
	host        = flag.String("host", "127.0.0.1", "")
	meterName   = flag.String("meter-name", "cloud-carbon", "")
)

func main() {
	start := time.Now()

	// setup
	ctx := context.Background()
	klog.InitFlags(nil)
	flag.Parse()

	// termination Handeling
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	// setup Server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", *host, *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// setup metrics
	metricshandler := metrics.New(*meterName)

	err := metricshandler.Setup(ctx)
	if err != nil {
		klog.Fatalf("failed initilizing metric handler: %v", err)
	}

	// register metrics endpoint
	http.Handle(*metricsPath, promhttp.Handler())

	// register health endpoint
	http.HandleFunc("/healthz", healthProbe)

	// run server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			klog.Fatalf("%s", fmt.Sprintf("failed to bind on %s: %v", *port, err))
		}
	}()

	klog.Infof("started in %v\n%v", time.Since(start), startUpLog)

	// wait for termination signals
	<-termChan

	// any Code to Gracefully Shutdown should be done here
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		klog.Fatal("Graceful Shutdown Failed")
	}

	klog.Info("Shutting Down Gracefully")
}
