package buffer

import (
	"context"
	"time"
)

type FlusherFunc = func(ctx context.Context, data Data) error

// Flush tries to flush current buffer data
// Block execution until flusher func succeeded or buffer is closed
func (b *Buffer) Flush() {
	flushData := b.swapData()
	if flushData.Empty() {
		return
	}

	for {
		if err := b.flushWithTimeout(flushData); err == nil {
			if b.throttling() {
				b.throttlingFlag.Store(false)
			}
			return
		}

		if b.closed() {
			return
		}

		b.throttlingFlag.Store(true)
		time.Sleep(b.opts.flushInterval())
	}
}

func (b *Buffer) swapData() Data {
	b.mt.Lock()
	defer b.mt.Unlock()

	oldData := b.data
	b.data = b.opts.dataConstructor()
	return oldData
}

func (b *Buffer) flushWithTimeout(data Data) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.opts.flushTimeout())
	defer cancel()

	return b.opts.flusherFunc(ctx, data)
}
