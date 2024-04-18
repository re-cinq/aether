package v1

import (
	"encoding/json"
	"errors"
)

// Provider where the resource consumption data is collected from
type Provider string

// ErrParsingProvider Error parsing the Provider
var ErrParsingProvider = errors.New("unsupported Provider")

// Provider constants
const (

	// Amazon web services API
	AWS Provider = awsString

	// Azure cloud API
	Azure Provider = azureString

	// Google cloud platform API
	GCP Provider = gcpString

	// Prometheus API for baremetal and kubernetes support
	Prometheus Provider = prometheusString

	// custom third party providers
	Custom Provider = customString

	// Constant string definitions
	awsString        = "aws"
	azureString      = "azure"
	gcpString        = "gcp"
	prometheusString = "prometheus"
	customString     = "custom"
)

// Providers Lookup map for listing all the supported providers
// as well as deserializing them
var Providers = map[string]Provider{
	awsString:        AWS,
	azureString:      Azure,
	gcpString:        GCP,
	prometheusString: Prometheus,
	customString:     Custom,
}

// Return the provider as string
func (p Provider) String() string {
	return string(p)
}

// Custom deserialization for Provider
func (p *Provider) UnmarshalJSON(data []byte) error {
	var value string

	// Unmarshall the bytes
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	// Make sure the unmarshalled string value exists
	if provider, ok := Providers[value]; !ok {
		return ErrParsingProvider
	} else {
		*p = provider
	}

	return nil
}
