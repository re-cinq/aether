package bus

type worker struct {

	// The worker own channel
	data chan Event

	// Waiting for the worker to conclude its tasks
	closing chan bool
}

func newWorker() worker {
	return worker{
		data:    make(chan Event, 1),
		closing: make(chan bool, 1),
	}
}
