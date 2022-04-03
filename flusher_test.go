package buffer_test

import (
	"github.com/quer3q/go-buffer"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestFlusher(t *testing.T) {
	t.Run("flush buffer when it's full", func(t *testing.T) {
		bufferSize := 1
		flusher := NewMockFlusher()
		b, err := buffer.New(
			buffer.WithSize(bufferSize),
			buffer.WithDataConstructor(NewMockDataConstructor()),
			buffer.WithFlusherFunc(flusher.FlusherFunc()),
			buffer.WithFlushInterval(func() time.Duration {
				return time.Minute
			}),
			buffer.WithFlushTimeout(func() time.Duration {
				return time.Minute
			}),
		)
		require.NoError(t, err)
		require.NotNil(t, b)

		require.NoError(t, b.Push(1))
		require.NoError(t, b.Push(2))
		require.NoError(t, b.Push(3))
		require.NoError(t, b.Push(4))
		require.NoError(t, b.Push(5))

		time.Sleep(time.Second * 1)
		require.ElementsMatch(t, flusher.Get(), []int{1, 2, 3, 4, 5})
	})
	t.Run("flush buffer using flush interval", func(t *testing.T) {
		bufferSize := 100
		flusher := NewMockFlusher()
		b, err := buffer.New(
			buffer.WithSize(bufferSize),
			buffer.WithDataConstructor(NewMockDataConstructor()),
			buffer.WithFlusherFunc(flusher.FlusherFunc()),
			buffer.WithFlushInterval(func() time.Duration {
				return time.Millisecond // Immediately flush
			}),
			buffer.WithFlushTimeout(func() time.Duration {
				return time.Minute
			}),
		)
		require.NoError(t, err)
		require.NotNil(t, b)

		require.NoError(t, b.Push(1))
		require.NoError(t, b.Push(2))
		require.NoError(t, b.Push(3))
		require.NoError(t, b.Push(4))
		require.NoError(t, b.Push(5))

		time.Sleep(time.Second)
		require.ElementsMatch(t, flusher.Get(), []int{1, 2, 3, 4, 5})
	})
}

func BenchmarkFlusher(b *testing.B) {
	flusher := NewMockFlusher()
	buff, err := buffer.New(
		buffer.WithSize(1), // small size
		buffer.WithDataConstructor(NewMockDataConstructor()),
		buffer.WithFlusherFunc(flusher.FlusherFunc()),
		buffer.WithFlushInterval(func() time.Duration {
			return time.Millisecond // Immediately flush
		}),
		buffer.WithFlushTimeout(func() time.Duration {
			return time.Minute
		}),
	)
	require.NoError(b, err)
	require.NotNil(b, buff)

	for i := 0; i < b.N; i++ {
		require.NoError(b, buff.Push(b))
	}
}
