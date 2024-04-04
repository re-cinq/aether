package amazon

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	v1 "github.com/re-cinq/aether/pkg/types/v1"
	"github.com/stretchr/testify/assert"
)

// Test the updateInstancesMap logic. The state of the
// instancesMap is shared among all the runs of the test
func TestUpdateInstancesMap(t *testing.T) {
	c := Client{
		instancesMap: make(map[string]*v1.Instance),
	}

	// A test up of the []types.Reservation type that contains
	// EC2 instance information
	res := []types.Reservation{
		{
			Instances: []types.Instance{
				{
					InstanceId:        aws.String("foo123"),
					InstanceType:      types.InstanceTypeA1Medium,
					InstanceLifecycle: types.InstanceLifecycleTypeScheduled,
					CpuOptions: &types.CpuOptions{
						CoreCount:      aws.Int32(4),
						ThreadsPerCore: aws.Int32(2),
					},
					State: &types.InstanceState{
						Name: types.InstanceStateNameRunning,
						Code: aws.Int32(16),
					},
				},
				{
					InstanceId:        aws.String("bar456"),
					InstanceType:      types.InstanceTypeT3Micro,
					InstanceLifecycle: types.InstanceLifecycleTypeScheduled,
					CpuOptions: &types.CpuOptions{
						CoreCount:      aws.Int32(4),
						ThreadsPerCore: aws.Int32(2),
					},
					State: &types.InstanceState{
						Name: types.InstanceStateNameTerminated,
						Code: aws.Int32(48),
					},
				},
			},
		},
	}

	t.Run("Running instance added to instancesMap", func(t *testing.T) {
		c.updateInstancesMap("fakeRegion", res)

		// check that the running instance is added to instancesMap
		_, exists := c.instancesMap["fakeRegion-AWS/EC2-foo123"]
		assert.True(t, exists)
	})

	t.Run("Terminated instance not added to instancesMap", func(t *testing.T) {
		c.updateInstancesMap("fakeRegion", res)

		// check that the terminated instance is not added to instancesMap
		_, exists := c.instancesMap["fakeRegion-AWS/EC2-bar456"]
		assert.False(t, exists)
	})

	t.Run("Instance in map changed to terminated state, remove from map", func(t *testing.T) {
		// check that the running instance still exists in the map
		_, exists := c.instancesMap["fakeRegion-AWS/EC2-foo123"]
		assert.True(t, exists)

		// modify the state of that instance to be "stopping"
		res[0].Instances[0].State = &types.InstanceState{
			Name: types.InstanceStateNameStopping,
			Code: aws.Int32(64),
		}

		c.updateInstancesMap("fakeRegion", res)

		// check that the instance was deleted from the instancesMap
		_, exists = c.instancesMap["fakeRegion-AWS/EC2-foo123"]
		assert.False(t, exists)
	})
}
