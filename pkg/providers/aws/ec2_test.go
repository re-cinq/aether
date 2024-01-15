package amazon

// TODO write unit tests that don't make API calls
//  func TestEc2InstanceListing(t *testing.T) {
//	  ctx := context.Background()
//	  // Pass an empty provider config so that it loads the default credentials
//	  cfg, err := buildAWSConfig(ctx, &config.Account{}, nil)
//	  assert.NotNil(t, cfg)
//	  assert.Nil(t, err)
//
//	  // Init the ec2 client
//	  ec2Client := NewEC2Client(&cfg)
//	  assert.NotNil(t, ec2Client)
//
//	  err = ec2Client.Refresh(ctx, "eu-north-1")
//	  assert.Nil(t, err)
//  }
