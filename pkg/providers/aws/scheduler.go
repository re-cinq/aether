package amazon

import (
	"time"

	"github.com/re-cinq/cloud-carbon/pkg/bus"
	"github.com/re-cinq/cloud-carbon/pkg/config"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	"k8s.io/klog/v2"
)

type awsScheduler struct {

	// Ticker
	ticker *time.Ticker

	// Signal we are done and shutting down
	done chan bool

	// Regions to scrape
	regions []string

	// Event bus
	eventBus *bus.EventBus

	// AWS Client
	// awsClient *AWSClient

	// Ec2 client
	ec2Client *ec2Client

	// Cloud watch client
	cloudwatchClient *cloudWatchClient
}

// Return the scheduler interface
func NewScheduler(eventBus *bus.EventBus) v1.Scheduler {

	// Init the ticket
	ticker := time.NewTicker(5 * time.Minute)

	// Get the list of regions
	regions := config.AppConfig().Providers[v1.Aws].Regions

	// Init the AWS Client
	awsClient, err := NewAWSClient()
	if err != nil {
		klog.Errorf("failed to initialise AWS provider %s", err)
		return nil
	}

	// Init the ec2 client
	ec2Client := NewEc2Client(awsClient.Config())
	if ec2Client == nil {
		klog.Fatal("Could not initialize EC2 client")
	}

	// Build the initial cache of instances
	for _, region := range regions {
		ec2Client.Refresh(region)
	}

	// Init the cloudwatch client
	cloudwatchClient := NewCloudWatchClient(awsClient.Config())
	if cloudwatchClient == nil {
		klog.Fatal("Could not initialize CloudWatch client")
	}

	return &awsScheduler{
		ticker:           ticker,
		done:             make(chan bool),
		eventBus:         eventBus,
		cloudwatchClient: cloudwatchClient,
		ec2Client:        ec2Client,
	}
}

func (s *awsScheduler) process() {

	for _, region := range s.regions {

		instances := s.cloudwatchClient.GetEc2Metrics(region, s.ec2Client.Cache())

		for _, instance := range instances {

			// Publish the metrics
			s.eventBus.Publish(v1.MetricsCollected{
				Instance: instance,
			})
		}

	}

}

func (s *awsScheduler) Schedule() {

	// Do the first call
	s.process()

	go func() {
		for {
			select {
			case <-s.done:
				return
			case <-s.ticker.C:
				s.process()
			}
		}
	}()

}

func (s *awsScheduler) Cancel() {

	// We are done
	s.done <- true

	// Stop the ticker
	s.ticker.Stop()

}
