package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/re-cinq/cloud-carbon/pkg/config"
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
		slog.Error("Could not create API router")
		os.Exit(1)
	}

	// Print to the user where the API is listening
	slog.Info("server started", "port", api.hostPort)

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
		slog.Error("failed to listen on port", "port", api.hostPort, "error", err)
		os.Exit(1)
	}
}

// Stop the api server
func (api *API) Stop() {
	slog.Info("Shutting down the API server")

	// Create a timeout for shutting down
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*15)

	// Release the context
	defer cancel()

	// Shutdown the server
	if err := api.server.Shutdown(ctxTimeout); err != nil {
		slog.Error("failed to gracefully shutdown the API", "error", err)
	}
}
