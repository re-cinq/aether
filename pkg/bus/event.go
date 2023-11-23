package bus

type Event interface {
	// Topic Returns the topic the event belongs to
	Topic() Topic

	// Id Returns the event unique Id which is used to distribute
	// the processing of the event itself across different workers
	Id() string
}
