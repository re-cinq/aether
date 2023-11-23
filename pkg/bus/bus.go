package bus

import (
	"hash/fnv"
	"sync"
	"time"

	"k8s.io/klog/v2"
)

type EventBus struct {
	// The size of the main queue
	queueSize int

	// The amount of parallel workers
	workerPoolSize int

	// Storing of the worker references
	workers map[uint32]worker

	// Storing of the subscribers to specific topics
	subscribers map[Topic][]EventHandler

	// The main queue
	queue chan Event

	// Locking for updating the various fields
	lock sync.RWMutex

	// Signaling whether the bus is shutting down, so no more messages will be accepted
	shutdown bool

	// Wait for the bus to shutdown gracefully
	done chan bool
}

// NewEventBus instance
func NewEventBus(queueSize int, workerPoolSize int) *EventBus {

	// Init the event bus
	bus := EventBus{
		queueSize:      queueSize,
		workerPoolSize: workerPoolSize,
		workers:        make(map[uint32]worker, workerPoolSize),
		subscribers:    make(map[Topic][]EventHandler),
		queue:          make(chan Event, queueSize),
		lock:           sync.RWMutex{},
		shutdown:       false,
		done:           make(chan bool),
	}

	// Init the channels and start the workers
	for i := 0; i < workerPoolSize; i++ {
		// init the channel
		bus.workers[uint32(i)] = newWorker()

		// Start the worker
		bus.startWorker(uint32(i))
	}

	return &bus
}

func (bus *EventBus) Subscribe(topic Topic, subscriber EventHandler) {
	bus.lock.Lock()
	if topicSubscribers, ok := bus.subscribers[topic]; ok {
		topicSubscribers = append(topicSubscribers, subscriber)
		bus.subscribers[topic] = topicSubscribers
	} else {
		bus.subscribers[topic] = []EventHandler{subscriber}
	}
	bus.lock.Unlock()
}

func (bus *EventBus) Start() {
	go bus.process()
}

func (bus *EventBus) Stop() {

	// shutdown
	bus.lock.Lock()
	bus.shutdown = true
	bus.lock.Unlock()

	// wait for the shutdown to be done
	<-bus.done

	klog.Info("event bus shutdown completed")
}

func (bus *EventBus) Publish(event Event) {
	var shuttingDown bool

	bus.lock.RLock()
	shuttingDown = bus.shutdown
	bus.lock.RUnlock()

	if !shuttingDown {
		bus.queue <- event
	}

}

func (bus *EventBus) process() {

	// Process all the events sequentially
	for {
		// shutdown the worker channels gracefully
		var shuttingDown bool

		bus.lock.RLock()
		shuttingDown = bus.shutdown
		bus.lock.RUnlock()

		if shuttingDown {
			break
		}

		select {
		// Get the event
		case event, more := <-bus.queue:
			if more {
				// Calculate the worker ID from the topic name
				workerId := bus.getWorkerId(event.Id())

				// Get the worker channel
				bus.lock.RLock()
				worker := bus.workers[workerId]
				bus.lock.RUnlock()

				// Send the event for processing
				worker.data <- event
			}
			// timeout otherwise to allow the shutdown procedure
		case <-time.After(5 * time.Second):
		}
	}

	// Shutdown the workers
	for _, workerChan := range bus.workers {
		// close it
		close(workerChan.data)

		// wait for the worker to terminate its tasks
		<-workerChan.closing
	}

	bus.done <- true

}

// Calculate the worker ID from the event topic
func (bus *EventBus) getWorkerId(eventId string) uint32 {

	eventHash := fnv.New32a()

	// Write the data
	eventHash.Write([]byte(eventId))

	// Calculate the id
	id := eventHash.Sum32()

	// Modulo the pool size
	return id % uint32(bus.workerPoolSize)

}

// Worker which processes the event
func (bus *EventBus) startWorker(id uint32) {

	// Get the channel
	workerChan := bus.workers[id]

	go func() {

		// Loop through it or wait for a task
		for {
			// Get the event
			event, more := <-workerChan.data

			// If we do have data
			if more {
				// load all the subscribers
				if subscribers, ok := bus.subscribers[event.Topic()]; ok {
					// Sync group
					var wg sync.WaitGroup

					// Loop through all subscribers
					for _, subscriber := range subscribers {
						// add to sync group
						wg.Add(1)
						go func(sub EventHandler) {
							// apply
							sub.Apply(event)

							// done
							wg.Done()
						}(subscriber)
					}

					// wait for all the event handlers to finish
					wg.Wait()

				}
			} else {
				// Channel was closed, no more data to be processed
				// Send the signal back that we are done here
				workerChan.closing <- true

				// Exit the for loop
				return
			}
		}
	}()

}
