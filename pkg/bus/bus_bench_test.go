package bus

import (
	"context"
	"sync"
	"testing"
)

// Define a mock EventHandler for testing
type mockHandler struct {
	mu    sync.Mutex
	count int
}

func (m *mockHandler) Handle(ctx context.Context, e *Event) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.count++
}

func (m *mockHandler) Stop(ctx context.Context) {
	// Do nothing
}

func BenchmarkBus(b *testing.B) {
	// Create a context for the benchmark
	ctx := context.Background()
	// Create a new bus instance
	bus := New(WithBufferSize(100), WithWorkers(10))

	var handlers []*mockHandler
	// Subscribe some mock handlers to the bus
	for i := 0; i < 1000; i++ {
		m := &mockHandler{}
		bus.Subscribe(EventType(i%10), m)

		handlers = append(handlers, m)
	}

	// Start the bus
	bus.Start(ctx)

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Publish an event to the bus
		event := &Event{
			Type: EventType(i % 10),
			Data: i,
		}
		err := bus.Publish(event)
		if err != nil {
			b.Fatalf("failed on publish: %v", err)
		}
	}

	// Stop the bus
	bus.Stop(ctx)

	expectedCount := b.N / 10
	for _, h := range handlers {
		// some of the handlers will receive either one less
		// or one more due to rounding
		// so this is an acceptable range
		lowerBound := expectedCount - 1
		upperBound := expectedCount + 1
		if h.count < lowerBound || h.count > upperBound {
			b.Fatalf("handler received %d events, expected approximately %d", h.count, expectedCount)
		}
	}
}
