package amazon

import (
	"testing"

	"github.com/re-cinq/cloud-carbon/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestCloudWatchMetrics(t *testing.T) {

	region := "eu-north-1"

	// Pass an empty provider config so that it loads the default credentials
	cfg, err := buildAWSConfig(config.Provider{})
	assert.NotNil(t, cfg)
	assert.Nil(t, err)

	// Init the ec2 client
	ec2Client := NewEc2Client(cfg)
	ec2Client.refresh(region)
	assert.NotNil(t, ec2Client)

	// Init the cloudwatch client
	client := NewCloudWatchClient(cfg)
	assert.NotNil(t, client)

	client.getEc2Metrics(region, ec2Client.cache)

	// assert.True(t, false)

}
