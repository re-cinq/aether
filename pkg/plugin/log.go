package plugin

import (
	"github.com/hashicorp/go-hclog"
	"github.com/re-cinq/aether/pkg/config"
)

// getLogLevel is a helper function to get the log level from the config
func getLogLevel() hclog.Level {
	level := config.AppConfig().LogLevel
	switch level {
	case "debug":
		return hclog.Debug
	case "info":
		return hclog.Info
	case "warn":
		return hclog.Warn
	case "error":
		return hclog.Error
	default:
		return hclog.Info
	}
}
