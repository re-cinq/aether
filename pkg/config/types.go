package config

import (
	"time"

	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
)

// ApplicationConfig is the full app config
type ApplicationConfig struct {
	APIConfig APIConfig                `mapstructure:"api"`
	Proxy     ProxyConfig              `mapstructure:"proxy"`
	Providers map[v1.Provider]Provider `mapstructure:"providers"`
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

// Defines the general configuration for a provider
type Provider struct {
	// The regions we should scrape the data for
	Regions []string `mapstructure:"regions"`

	// AWS Specific:
	// Cloudwatch namespaces
	// A namespace is a container for CloudWatch metrics.
	// Metrics in different namespaces are isolated from each other,
	// so that metrics from different applications are not mistakenly aggregated into the same statistics.
	// For example, Amazon EC2 uses the AWS/EC2 namespace.
	// For the list of AWS namespaces, see AWS services that publish CloudWatch metrics.
	// https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/aws-services-cloudwatch-metrics.html
	Namespaces []string `mapstructure:"namespaces"`

	// The location from where to load the credentials
	Credentials ProviderConfig `mapstructure:"credentials"`

	// The location from where to load the additional configuration
	Config ProviderConfig `mapstructure:"config"`

	// The SDK Http Client transport configuration
	Transport TransportConfig `mapstructure:"transport"`
}

type ProviderConfig struct {
	// AWS: which profile to use
	Profile string `mapstructure:"profile"`

	// Where the file can be located
	FilePaths []string `mapstructure:"filePaths"`
}

// Whether the provider config has some values set
func (pc ProviderConfig) IsPresent() bool {
	return pc.Profile != "" || len(pc.FilePaths) > 0
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
