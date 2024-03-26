package bus

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type testHandler struct {
	received *Event
	wg       *sync.WaitGroup // Synchronization mechanism
}

func (h *testHandler) Handle(ctx context.Context, e *Event) {
	h.received = e
	h.wg.Done()
}

func (h *testHandler) Stop(ctx context.Context) {
	h.wg.Wait() // Wait until all handles are done
}

func TestBus(t *testing.T) {
	assert := require.New(t)

	ctx := context.Background()
	var topic EventType = 1

	h1 := &testHandler{wg: &sync.WaitGroup{}}
	h2 := &testHandler{wg: &sync.WaitGroup{}}

	h1.wg.Add(1)
	h2.wg.Add(1)

	b := New(WithWorkers(1), WithBufferSize(10))

	b.Subscribe(topic, h1)
	b.Subscribe(topic, h2)

	b.Start(ctx)

	e := &Event{
		Type: topic,
		Data: "test",
	}

	err := b.Publish(e)
	assert.NoError(err)

	h1.wg.Wait()
	h2.wg.Wait()

	assert.Equal(h1.received, e)
	assert.Equal(h2.received, e)

	e2 := &Event{
		Type: topic,
		Data: "test",
	}

	h1.wg.Add(1)
	h2.wg.Add(1)
	err = b.Publish(e)
	assert.NoError(err)

	h1.wg.Wait()
	h2.wg.Wait()

	assert.Equal(h1.received, e2)
	assert.Equal(h2.received, e2)

	b.Stop(ctx)
}
