package buffer_test

import (
	"context"
	"github.com/quer3q/go-buffer"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestBuffer(t *testing.T) {
	t.Run("can't init empty buffer", func(t *testing.T) {
		b, err := buffer.New()
		require.Error(t, err)
		require.Nil(t, b)
	})
}

type (
	MockDataConstructor struct {
		values []interface{}
	}

	MockFlusher struct {
		values []interface{}
		calls  int
		m      *sync.Mutex
	}
)

func (f *MockFlusher) FlusherFunc() buffer.FlusherFunc {
	return func(ctx context.Context, data buffer.Data) error {
		f.m.Lock()
		defer f.m.Unlock()
		v, _ := data.(*MockDataConstructor)
		f.values = append(f.values, v.Get()...)
		f.calls++
		return nil
	}
}

func (f *MockFlusher) Get() []interface{} {
	f.m.Lock()
	defer f.m.Unlock()
	return f.values
}

func (m *MockDataConstructor) Push(v interface{}) int {
	m.values = append(m.values, v)
	return len(m.values)
}

func (m *MockDataConstructor) Empty() bool {
	return len(m.values) == 0
}

func (m *MockDataConstructor) Get() []interface{} {
	return m.values
}

func NewMockDataConstructor() buffer.DataConstructor {
	return func() buffer.Data {
		return &MockDataConstructor{}
	}
}

func NewMockFlusher() *MockFlusher {
	return &MockFlusher{m: &sync.Mutex{}}
}
