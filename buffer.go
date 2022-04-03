package buffer

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

var (
	errBufferClosed     = errors.New("buffer is closed")
	errBufferThrottling = errors.New("buffer is throttling")
)

type Buffer struct {
	data           Data
	mt             *sync.Mutex
	opts           *Options
	closedFlag     atomic.Value
	throttlingFlag atomic.Value
	closedCh       chan struct{}
	mustFlushCh    chan struct{}
}

func New(opts ...Option) (*Buffer, error) {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	if err := validateOptions(options); err != nil {
		return nil, err
	}

	var closed, throttled atomic.Value
	closed.Store(false)
	throttled.Store(false)

	b := &Buffer{
		data:           options.dataConstructor(),
		mt:             &sync.Mutex{},
		opts:           options,
		closedCh:       make(chan struct{}),
		mustFlushCh:    make(chan struct{}),
		closedFlag:     closed,
		throttlingFlag: throttled,
	}

	go b.consume()
	return b, nil
}

func (b *Buffer) Push(v interface{}) error {
	if err := b.checkBuffer(); err != nil {
		return errors.Wrap(err, "buffer is not working")
	}

	b.mt.Lock()
	defer b.mt.Unlock()

	if total := b.data.Push(v); total >= b.opts.size {
		go func() {
			b.mustFlushCh <- struct{}{}
		}()
	}

	return nil
}

func (b *Buffer) Close() error {
	if b.closed() {
		return errBufferClosed
	}

	b.closedFlag.Store(true)
	b.closedCh <- struct{}{}
	return nil
}

func (b *Buffer) consume() {
	ticker := time.NewTicker(b.opts.flushInterval())

	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			b.Flush()
			ticker = time.NewTicker(b.opts.flushInterval())
		case <-b.mustFlushCh:
			if !b.throttling() {
				b.Flush()
			}
		case <-b.closedCh:
			b.Flush()
			return
		}
	}
}

func (b *Buffer) closed() bool {
	return b.closedFlag.Load().(bool)
}

func (b *Buffer) throttling() bool {
	return b.throttlingFlag.Load().(bool)
}

func (b *Buffer) checkBuffer() error {
	if b.closed() {
		return errBufferClosed
	}
	if b.throttling() {
		return errBufferThrottling
	}

	return nil
}
