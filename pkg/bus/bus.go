package bus

import (
	"context"
	"errors"
	"sync"
)

var (
	// This is the default buffer size of the queue
	defaultbuffer = 512

	// the default amount of workers that the bus will start with
	defaultworkers = 10
)

// EventType is the identifier for the type of event you are dealing with
type EventType int

// Bus is the main structure used to keep track of all subscribers
type Bus struct {
	// holds event subscribers
	subs map[EventType][]EventHandler

	// The queue
	queue chan *Event

	// amount of workers to start
	workers int

	// termination signal
	shutdown bool

	// syncing functionality
	wg sync.WaitGroup
	mu sync.RWMutex
}

// Event is the main type that gets sent thr9ough the bus
// it is up to the callers to validate the data when received
type Event struct {
	Type EventType
	Data interface{}
}

// EventHandler is an interface that callers can use
// to make use of the bus
type EventHandler interface {
	// This function will be called when an event is received
	// that the handler is subscribed to
	Handle(ctx context.Context, e *Event)
	// On Shutfown of the bus all Stop functyions will be called
	// NOTE: Stop should be idempotent as it will be called once per
	// subscribed event. so if a handler is subscribed to more than one event it
	// will be called multiple times
	Stop(ctx context.Context)
}

type option func(*Bus)

// WithBufferSize is used to set the buffer size of each queue
// each subscriber gets a channel that its listening on
// without a high enough buffer size the channel will stall
func WithBufferSize(s int) option {
	return func(b *Bus) {
		b.queue = make(chan *Event, s)
	}
}

// WithWorkers sets the amount of workers that will be used
func WithWorkers(w int) option {
	return func(b *Bus) {
		b.workers = w
	}
}

// New instanciates a new instance of Bus
func New(opts ...option) *Bus {
	// sets the defaults
	b := &Bus{
		subs:     make(map[EventType][]EventHandler),
		queue:    make(chan *Event, defaultbuffer),
		shutdown: false,
		workers:  defaultworkers,
	}

	// sets any options
	for _, o := range opts {
		o(b)
	}

	return b
}

// Subscribe adds a handler to listen for a specific event type
func (b *Bus) Subscribe(e EventType, h EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.subs[e]; !ok {
		b.subs[e] = []EventHandler{}
	}
	b.subs[e] = append(b.subs[e], h)
}

// Publish a new event
func (b *Bus) Publish(e *Event) error {
	// if shutdown we no longer allow publishing
	if b.shutdown {
		return errors.New("bus has shutdown")
	}

	b.queue <- e

	return nil
}

// Start Bus by starting the workers
func (b *Bus) Start(ctx context.Context) {
	for i := 0; i < b.workers; i++ {
		b.wg.Add(1)
		go b.worker(ctx)
	}
}

// worker processes events from the queue
// on shutdown signal it will return
// we do allow the ability to pass trhough a context to force the worker to
// shutdown
func (b *Bus) worker(ctx context.Context) {
	defer b.wg.Done()
	for {
		select {
		case event, ok := <-b.queue:
			if !ok {
				// The queue has been closed, so the worker exits
				return
			}

			b.mu.RLock()
			subs := b.subs[event.Type]
			b.mu.RUnlock()

			for _, sub := range subs {
				// We make a copy of the event to reduce
				// the likelihood of race conditions
				e := *event

				sub.Handle(ctx, &e)
			}
		case <-ctx.Done():
			// if context is canceled
			// we should stop the worker
			return
		}
	}
}

// Stop the Bus gracefully by waiting for workers to exit
func (b *Bus) Stop(ctx context.Context) {
	b.mu.Lock()
	b.shutdown = true
	close(b.queue)
	b.mu.Unlock()
	b.wg.Wait()
}
