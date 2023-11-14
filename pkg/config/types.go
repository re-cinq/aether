package config

import "time"

// ApplicationConfig is the full app config
type ApplicationConfig struct {
	APIConfig APIConfig   `mapstructure:"api"`
	Proxy     ProxyConfig `mapstructure:"proxy"`
	// Providers     ProviderConfigs `mapstructure:"providers"`
	// ServiceConfig ServiceConfig   `mapstructure:"service"`
	// DBConfig      DBConfig        `mapstructure:"pg"`
}

// Defines the configuration for the API
type APIConfig struct {
	// The address to listen to
	Address string `mapstructure:"address"`

	// The port to listen to
	Port string `mapstructure:"port"`
}

// Generic proxy configuration
type ProxyConfig struct {
	HttpProxy  string `mapstructure:"httpProxy"`
	HttpsProxy string `mapstructure:"httpsProxy"`
	NoProxy    string `mapstructure:"noProxy"`
}

// Defines various Http Client transport settings
type TransportConfig struct {
	// Proxy
	Proxy ProxyConfig `mapstructure:"proxy"`

	// Idle connection timeout
	IdleConnTimeout time.Duration `mapstructure:"idleConnTimeout"`

	// Max idle connections
	MaxIdleConns int `mapstructure:"maxIdleConns"`

	// maximum number of idle (keep-alive) connections per-host.
	MaxIdleConnsPerHost int `mapstructure:"maxIdleConnsPerHost"`

	// Timeout while reading headers
	ResponseHeaderTimeout time.Duration `mapstructure:"responseHeaderTimeout"`

	// maximum amount of time waiting for a TLS handshake to be completed
	TLSHandshakeTimeout time.Duration `mapstructure:"tlsHandshakeTimeout"`
}
