package amazon

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/re-cinq/cloud-carbon/pkg/config"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"

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

// // Calls the AWS API to retrieve the metrics
// func (awsClient *AWSClient) getMetrics(namespace string, interval time.Duration) {

// 	// Load the current config
// 	currentConfig := config.AppConfig().Providers[v1.Aws]

// 	// Build the client
// 	client := cloudwatch.NewFromConfig(awsClient.cfg)

// 	input := &cloudwatch.GetMetricDataInput{
// 		EndTime:   aws.Time(time.Unix(time.Now().Unix(), 0)),
// 		StartTime: aws.Time(time.Unix(time.Now().Add(-interval).Unix(), 0)),
// 		MetricDataQueries: []types.MetricDataQuery{
// 			types.MetricDataQuery{
// 				Id: aws.String(*id),
// 				MetricStat: &types.MetricStat{
// 					Metric: &types.Metric{
// 						Namespace:  aws.String(*namespace),
// 						MetricName: aws.String(*metricName),
// 						Dimensions: []types.Dimension{
// 							types.Dimension{
// 								Name:  aws.String(*dimensionName),
// 								Value: aws.String(*dimensionValue),
// 							},
// 						},
// 					},
// 					Period: aws.Int32(int32(*period)),
// 					Stat:   aws.String(*stat),
// 				},
// 			},
// 		},
// 	}

// 	result, err := GetMetrics(context.TODO(), client, input)
// 	if err != nil {
// 		fmt.Println("Could not fetch metric data")
// 	}

// 	fmt.Println("Metric Data:", result)
// }

// NewAWSClient creates a new instance of the AWSClient
// It allows to pass:
//   - configFile: the location of the config file to load. If empty the default
//     location of the credentials file (~/.aws/config) is used
//   - profile: the name of the profile to use to load the credentials
//     if empty the default credentials will be used
func NewAWSClient() (*AWSClient, error) {

	currentConfig := config.AppConfig().Providers[v1.Aws]

	cfg, err := buildAWSConfig(currentConfig)

	if err != nil {
		klog.Errorf("failed to initialise AWS Client: %s", err)
		return nil, err
	}

	return &AWSClient{
		cfg: cfg,
	}, nil
}

// Helper function to builde the AWS config
func buildAWSConfig(currentConfig config.Provider) (aws.Config, error) {

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
	// .WithTransportOptions(func(tr *http.Transport) {
	// 	proxyURL, err := url.Parse("PROXY_URL")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	tr.Proxy = http.ProxyURL(proxyURL)
	// })

	loadExternalConfigs = append(loadExternalConfigs, awsConfig.WithHTTPClient(httpClient))

	// -------------------------------------------------------------------
	// Finally generate the config
	cfg, err = awsConfig.LoadDefaultConfig(context.TODO(), loadExternalConfigs...)

	return cfg, err

}
