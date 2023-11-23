package amazon

import (
	"testing"

	"github.com/re-cinq/cloud-carbon/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestEc2InstanceListing(t *testing.T) {

	// Pass an empty provider config so that it loads the default credentials
	cfg, err := buildAWSConfig(config.Account{}, nil)
	assert.NotNil(t, cfg)
	assert.Nil(t, err)

	// Init the ec2 client
	ec2Client := NewEc2Client(cfg)
	assert.NotNil(t, ec2Client)

	ec2Client.Refresh("eu-north-1")

	// assert.True(t, false)

}
