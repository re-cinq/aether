package amazon

import (
	"context"
	"errors"
	"fmt"

	"github.com/re-cinq/cloud-carbon/pkg/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

var (
	ErrLoadingAwsConfigFile = errors.New("failed to load AWS credentials config file")
)

// Config for AWS for calling API and setting up connections
type Client struct {
	config     *aws.Config
	ec2        *ec2Client
	cloudwatch *cloudWatchClient
}

// CWGetMetricDataAPI defines the interface for the GetMetricData function
type CWGetMetricDataAPI interface {
	GetMetricData(ctx context.Context, params *cloudwatch.GetMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error)
}

// GetMetrics Fetches the cloudwatch metrics for your provided input in the given time-frame
func (c *Client) GetMetrics(ctx context.Context, api CWGetMetricDataAPI, input *cloudwatch.GetMetricDataInput) (*cloudwatch.GetMetricDataOutput, error) {
	fmt.Printf("IS THIS ACTUALLY CALLED???????: %+v\n", input)
	return api.GetMetricData(ctx, input)
}

// NewAWSConfig creates a new AWS config for creating clients
// It allows to pass:
//   - configFile: the location of the config file to load. If empty the default
//     location of the credentials file (~/.aws/config) is used
//   - profile: the name of the profile to use to load the credentials
//     if empty the default credentials will be used
func NewAWSClient(ctx context.Context, currentConfig *config.Account, customTransportConfig *config.TransportConfig) (*Client, error) {
	cfg, err := buildAWSConfig(ctx, currentConfig, customTransportConfig)
	if err != nil {
		return nil, err
	}

	// Init the ec2 client
	ec2Client := NewEc2Client(&cfg)
	if ec2Client == nil {
		return nil, errors.New("Could not initialize EC2 client")
	}

	// Init the cloudwatch client
	cloudwatchClient := NewCloudWatchClient(&cfg)
	if cloudwatchClient == nil {
		return nil, errors.New("Could not initialize CloudWatch client")
	}

	return &Client{
		&cfg,
		ec2Client,
		cloudwatchClient,
	}, nil
}

// Helper function to builde the AWS config
func buildAWSConfig(ctx context.Context, currentConfig *config.Account, customTransportConfig *config.TransportConfig) (aws.Config, error) {
	//	var err error

	return awsConfig.LoadDefaultConfig(ctx)
	// -------------------------------------------------------------------

	// If the user did not pass the location of the config file to load, fall back
	// to the default location
	// Override the credentials and the config if necessary

	//loadExternalConfigs := []func(*awsConfig.LoadOptions) error{}
	//hasCredentials := len(currentConfig.Credentials.FilePaths) > 0
	//hasConfig := len(currentConfig.Config.FilePaths) > 0

	//// If we have credentials to override
	//if hasCredentials {
	//	loadExternalConfigs = append(loadExternalConfigs, awsConfig.WithSharedCredentialsFiles(currentConfig.Credentials.FilePaths))
	//}

	//// If there is a profile set
	//if currentConfig.Credentials.Profile != "" {
	//	loadExternalConfigs = append(loadExternalConfigs, awsConfig.WithSharedConfigProfile(currentConfig.Credentials.Profile))
	//}

	//// If we have configs to override
	//if hasConfig {
	//	loadExternalConfigs = append(loadExternalConfigs, awsConfig.WithSharedConfigFiles(currentConfig.Config.FilePaths))
	//}

	//// If there is a profile set
	//if currentConfig.Config.Profile != "" {
	//	loadExternalConfigs = append(loadExternalConfigs, awsConfig.WithSharedConfigProfile(currentConfig.Config.Profile))
	//}

	//// -------------------------------------------------------------------
	//// Http client
	//httpClient := awshttp.NewBuildableClient().WithDialerOptions(func(d *net.Dialer) {
	//	d.KeepAlive = -1
	//	d.Timeout = time.Millisecond * 500
	//})

	//// TODO: IS this correct?
	//// first parse HTTP Proxy url, then if HTTPS url exists, override that
	//// with HTTPS and set it in the custom transport
	//if customTransportConfig != nil {
	//	// Override the transport settings
	//	var proxyURL *url.URL
	//	if customTransportConfig.Proxy.HTTPProxy != "" {
	//		proxyURL, err = url.Parse(customTransportConfig.Proxy.HTTPProxy)
	//		if err != nil {
	//			klog.Fatalf("failed to parse config 'HTTPProxy' url")
	//		}
	//	}

	//	if customTransportConfig.Proxy.HTTPSProxy != "" {
	//		proxyURL, err = url.Parse(customTransportConfig.Proxy.HTTPSProxy)
	//		if err != nil {
	//			klog.Fatalf("failed to parse config 'HTTPSProxy' url")
	//		}
	//	}

	//	var customTransport *http.Transport
	//	if proxyURL != nil {
	//		customTransport.Proxy = http.ProxyURL(proxyURL)
	//	}

	//	// TODO: check all the additional transport settings and if different from the default override them

	//	// httpClient.WithTransportOptions(func(t *http.Transport) {
	//	// 	if customTransport.Proxy != nil {
	//	// 		t.Proxy = customTransport.Proxy
	//	// 	}

	//	// })
	//}

	//loadExternalConfigs = append(loadExternalConfigs, awsConfig.WithHTTPClient(httpClient))

	//// -------------------------------------------------------------------
	//// Finally generate the config
	//return awsConfig.LoadDefaultConfig(ctx, loadExternalConfigs...)
}
