package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/re-cinq/aether/pkg/config"
	"github.com/re-cinq/aether/pkg/log"
)

const readHeaderTimeout = 2 * time.Second

// ApiServer object
type API struct {
	*http.Server

	addr        string
	metricsPath string
}

// New returns an instance of a configured API
func New() *API {
	api := &API{
		metricsPath: config.AppConfig().APIConfig.MetricsPath,
		addr: fmt.Sprintf("%s:%s",
			config.AppConfig().APIConfig.Address,
			config.AppConfig().APIConfig.Port,
		),
	}

	api.setup()

	return api
}

// router configures the routes for the API server
func (a *API) router() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)

	// HealthCheck
	r.HandleFunc("/healthz", healthProbe).Methods("GET")

	// Prometheus exporter
	prometheus.MustRegister(version.NewCollector("cloud_carbon_exporter"))
	r.Handle(a.metricsPath, promhttp.Handler()).Methods("GET")

	return r
}

// setup configures the API server
func (a *API) setup() {
	// Init the HTTP server
	a.Server = &http.Server{
		Addr:              a.addr,
		Handler:           a.router(),
		ReadHeaderTimeout: readHeaderTimeout,
	}
}

// Start is a blocking event that starts the server on the configured port
func (a *API) Start(ctx context.Context) {
	logger := log.FromContext(ctx)
	logger.Info("server started", "address", a.addr)

	// Listen to it
	if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("failed to listen on address", "address", a.addr, "error", err)
		os.Exit(1)
	}
}

// Stop sends a termination signal to the api server
func (a *API) Stop(ctx context.Context) {
	logger := log.FromContext(ctx)
	logger.Info("shutting down the api server")

	// Shutdown the server
	if err := a.Server.Shutdown(ctx); err != nil {
		logger.Error("failed to gracefully shutdown the API", "error", err)
	}
}
