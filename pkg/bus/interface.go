package bus

type Bus interface {
	// Subscribe handler to topic
	Subscribe(topic Topic, subscriber EventHandler)

	// Start the bus
	Start()

	// Stop the bus
	Stop()

	// Publish the event to the bus
	Publish(event Event)
}
