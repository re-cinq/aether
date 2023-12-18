package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/re-cinq/cloud-carbon/pkg/config"
	"k8s.io/klog/v2"
)

// ApiServer object
type API struct {
	host        string
	port        string
	hostPort    string
	metricsPath string
	server      *http.Server
}

// NewAPIServer instance
func NewAPIServer() *API {
	// Local variable
	host := config.AppConfig().APIConfig.Address
	port := config.AppConfig().APIConfig.Port

	// Make sure the metrics path is set
	metricsPath := config.AppConfig().APIConfig.MetricsPath
	if metricsPath == "" {
		metricsPath = "/metrics"
	}

	// Return the API
	api := API{
		host:        host,
		port:        port,
		metricsPath: metricsPath,
		hostPort:    fmt.Sprintf("%s:%s", host, port),
	}

	// Initialize it
	api.init()

	// Return it
	return &api
}

func (api *API) createRouter() *httprouter.Router {
	// Configure the HTTP router
	router := httprouter.Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
	}

	// HealthCheck
	router.GET("/v1/health", healthProbe)

	// Prometheus exporter
	prometheus.MustRegister(version.NewCollector("cloud_carbon_exporter"))
	router.Handler("GET", api.metricsPath, promhttp.Handler())

	return &router
}

// Start the api server
func (api *API) init() {
	// Create the router and all the handlers
	router := api.createRouter()

	// Make sure we have got a router
	if router == nil {
		klog.Fatal("Could not create API router")
		return
	}

	// Print to the user where the API is listening
	klog.Infof("listening on: %s", api.hostPort)

	// Init the HTTP server
	api.server = &http.Server{
		Addr:              api.hostPort,
		Handler:           router,
		ReadHeaderTimeout: 2 * time.Second,
	}
}

func (api *API) Start() {
	// Listen to it
	if err := api.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		klog.Fatalf("failed to listen to %s %s", api.hostPort, err)
	}
}

// Stop the api server
func (api *API) Stop() {
	klog.Infof("Shutting down the API server...")

	// Create a timeout for shutting down
	ctxTimeout, cancel := context.WithTimeout(context.TODO(), time.Second*15)

	// Release the context
	defer cancel()

	// Shutdown the server
	if err := api.server.Shutdown(ctxTimeout); err != nil {
		klog.Errorf("failed to gracefully shutdown the API: %s", err)
	}
}
