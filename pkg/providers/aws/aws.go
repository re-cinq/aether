package amazon

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/re-cinq/cloud-carbon/pkg/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"k8s.io/klog/v2"
)

var (
	ErrLoadingAwsConfigFile = errors.New("failed to load AWS credentials config file")
)

// AWSClient is responsible for calling the AWS API
type AWSClient struct {
	// The config in case we have to re-establish a connection
	cfg aws.Config
}

// CWGetMetricDataAPI defines the interface for the GetMetricData function
type CWGetMetricDataAPI interface {
	GetMetricData(ctx context.Context, params *cloudwatch.GetMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error)
}

// GetMetrics Fetches the cloudwatch metrics for your provided input in the given time-frame
func (awsClient *AWSClient) GetMetrics(c context.Context, api CWGetMetricDataAPI, input *cloudwatch.GetMetricDataInput) (*cloudwatch.GetMetricDataOutput, error) {
	return api.GetMetricData(c, input)
}

// NewAWSClient creates a new instance of the AWSClient
// It allows to pass:
//   - configFile: the location of the config file to load. If empty the default
//     location of the credentials file (~/.aws/config) is used
//   - profile: the name of the profile to use to load the credentials
//     if empty the default credentials will be used
func NewAWSClient(currentConfig config.Account, customTransportConfig *config.TransportConfig) (*AWSClient, error) {

	cfg, err := buildAWSConfig(currentConfig, customTransportConfig)

	if err != nil {
		klog.Errorf("failed to initialize AWS Client: %s", err)
		return nil, err
	}

	return &AWSClient{
		cfg: cfg,
	}, nil
}

func (awsClient *AWSClient) Config() aws.Config {
	return awsClient.cfg
}

// Helper function to builde the AWS config
func buildAWSConfig(currentConfig config.Account, customTransportConfig *config.TransportConfig) (aws.Config, error) {

	// Define the variables to be populated based on the provider configuration
	// AWS config file
	var cfg aws.Config

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
				klog.Fatalf("failed to parse config 'HTTPProxy' url")
			}
		}

		if customTransportConfig.Proxy.HTTPSProxy != "" {
			proxyURL, err = url.Parse(customTransportConfig.Proxy.HTTPSProxy)
			if err != nil {
				klog.Fatalf("failed to parse config 'HTTPSProxy' url")
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
	cfg, err = awsConfig.LoadDefaultConfig(context.TODO(), loadExternalConfigs...)

	return cfg, err
}
