package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/re-cinq/cloud-carbon/pkg/config"
	amazon "github.com/re-cinq/cloud-carbon/pkg/providers/aws"
	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
	"k8s.io/klog/v2"
)

func main() {

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	// Init the config
	config.InitConfig()

	region := config.AppConfig().Providers[v1.Aws].Regions[0]

	awsClient, err := amazon.NewAWSClient()
	if err != nil {
		klog.Fatal(err)
	}

	// Init the ec2 client
	ec2Client := amazon.NewEc2Client(awsClient.Config())
	if ec2Client == nil {
		klog.Fatal("Could not initialize EC2 client")
	}
	ec2Client.Refresh(region)

	// Init the cloudwatch client
	client := amazon.NewCloudWatchClient(awsClient.Config())
	if client == nil {
		klog.Fatal("Could not initialize CloudWatch client")
	}

	client.GetEc2Metrics(region, ec2Client.Cache())

	// -----------------------------------------------------------

	scheduler := time.NewTicker(30 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-scheduler.C:
				ec2Client.Refresh(region)
				client.GetEc2Metrics(region, ec2Client.Cache())
			}
		}
	}()

	<-termChan
	done <- true

}
