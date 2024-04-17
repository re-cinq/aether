package amazon

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/re-cinq/aether/pkg/config"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

var (
	ErrLoadingAwsConfigFile = errors.New("failed to load AWS credentials config file")
)

// Client contains the AWS config and service clients
// and is used to access the API
type Client struct {
	// AWS specific config for auth and client creation
	cfg *aws.Config

	// service APIs
	ec2        *ec2.Client
	cloudwatch *cloudwatch.Client

	instancesMap map[string]*v1.Instance
}

// New creates a struct with the AWS config, EC2 Client, and CloudWatch Client
// It allows to pass:
//   - configFile: the location of the config file to load. If empty the default
//     location of the credentials file (~/.aws/config) is used
//   - profile: the name of the profile to use to load the credentials
//     if empty the default credentials will be used
//
// TODO: use options pattern
func New(ctx context.Context, currentConfig *config.Account, customTransportConfig *config.TransportConfig) (*Client, error) {
	cfg, err := buildAWSConfig(ctx, currentConfig, customTransportConfig)
	if err != nil {
		return nil, fmt.Errorf("error initializing AWS client: %s", err)
	}
	c := &Client{
		cfg:          cfg,
		instancesMap: make(map[string]*v1.Instance),
	}

	// TODO: fix options pattern
	c.ec2 = ec2.NewFromConfig(*cfg, func(o *ec2.Options) {})
	if c.ec2 == nil {
		return nil, fmt.Errorf("unable to initialize ec2 client")
	}

	// TODO: fix options pattern
	c.cloudwatch = cloudwatch.NewFromConfig(*cfg, func(o *cloudwatch.Options) {})
	if c.cloudwatch == nil {
		return nil, errors.New("error initializing CloudWatch client")
	}

	return c, nil
}

// Helper function to builde the AWS config
func buildAWSConfig(ctx context.Context, currentConfig *config.Account, customTransportConfig *config.TransportConfig) (*aws.Config, error) {
	// Error when loading the config file
	var err error

	// -------------------------------------------------------------------

	// If the user did not pass the location of the config file to load, fall back
	// to the default location
	// Override the credentials and the config if necessary

	loadExternalConfigs := []func(*awsConfig.LoadOptions) error{}
	hasCredentials := len(currentConfig.Credentials.FilePaths) > 0
	hasConfig := len(currentConfig.Config.FilePaths) > 0

	// If we have credentials to override
	if hasCredentials {
		loadExternalConfigs = append(loadExternalConfigs, awsConfig.WithSharedCredentialsFiles(currentConfig.Credentials.FilePaths))
	}

	// If there is a profile set
	if currentConfig.Credentials.Profile != "" {
		loadExternalConfigs = append(loadExternalConfigs, awsConfig.WithSharedConfigProfile(currentConfig.Credentials.Profile))
	}

	// If we have configs to override
	if hasConfig {
		loadExternalConfigs = append(loadExternalConfigs, awsConfig.WithSharedConfigFiles(currentConfig.Config.FilePaths))
	}

	// If there is a profile set
	if currentConfig.Config.Profile != "" {
		loadExternalConfigs = append(loadExternalConfigs, awsConfig.WithSharedConfigProfile(currentConfig.Config.Profile))
	}

	// -------------------------------------------------------------------
	// Http client
	httpClient := awshttp.NewBuildableClient().WithDialerOptions(func(d *net.Dialer) {
		d.KeepAlive = -1
		d.Timeout = time.Millisecond * 500
	})

	if customTransportConfig != nil {
		// Override the transport settings
		var proxyURL *url.URL
		if customTransportConfig.Proxy.HTTPProxy != "" {
			proxyURL, err = url.Parse(customTransportConfig.Proxy.HTTPProxy)
			if err != nil {
				return nil, err
			}
		}

		if customTransportConfig.Proxy.HTTPSProxy != "" {
			proxyURL, err = url.Parse(customTransportConfig.Proxy.HTTPSProxy)
			if err != nil {
				return nil, err
			}
		}

		var customTransport *http.Transport
		if proxyURL != nil {
			customTransport.Proxy = http.ProxyURL(proxyURL)
		}

		// TODO: check all the additional transport settings and if different from the default override them

		// httpClient.WithTransportOptions(func(t *http.Transport) {
		// 	if customTransport.Proxy != nil {
		// 		t.Proxy = customTransport.Proxy
		// 	}

		// })
	}

	loadExternalConfigs = append(loadExternalConfigs, awsConfig.WithHTTPClient(httpClient))

	// -------------------------------------------------------------------
	// Finally generate the config
	c, err := awsConfig.LoadDefaultConfig(ctx, loadExternalConfigs...)
	return &c, err
}
