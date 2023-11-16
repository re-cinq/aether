package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLabelsOperations(t *testing.T) {
	key := "name"
	value := "test"

	// Init an empty var
	labels := Labels{}

	// Add a label
	labels.Add(key, value)

	// Make sure it exists
	assert.True(t, labels.Exists(key))

	// Get the value back
	existingValue, exists := labels.Get(key)
	assert.True(t, exists)
	assert.Equal(t, value, existingValue)

	// Delete an entry
	labels.Delete(key)

	// Make sure the entry was deleted
	assert.False(t, labels.Exists(key))
}
