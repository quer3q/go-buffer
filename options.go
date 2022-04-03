package buffer

import (
	"time"

	"github.com/pkg/errors"
)

type (
	Options struct {
		size            int
		dataConstructor DataConstructor
		flushInterval   func() time.Duration
		flushTimeout    func() time.Duration
		flusherFunc     FlusherFunc
	}

	Option func(*Options)
)

func WithSize(size int) Option {
	return func(o *Options) {
		o.size = size
	}
}

func WithFlushTimeout(fn func() time.Duration) Option {
	return func(o *Options) {
		o.flushTimeout = fn
	}
}

func WithFlushInterval(fn func() time.Duration) Option {
	return func(o *Options) {
		o.flushInterval = fn
	}
}

func WithDataConstructor(fn DataConstructor) Option {
	return func(o *Options) {
		o.dataConstructor = fn
	}
}

func WithFlusherFunc(fn FlusherFunc) Option {
	return func(o *Options) {
		o.flusherFunc = fn
	}
}

func validateOptions(o *Options) error {
	if o.size <= 0 {
		return errors.New("max size can't be less or equal to 0")
	}
	if o.flusherFunc == nil {
		return errors.New("flusher func can't be empty")
	}
	if o.dataConstructor == nil {
		return errors.New("data constructor func can't be empty")
	}
	if o.flushInterval == nil || o.flushInterval() <= 0 {
		return errors.New("flush interval can't be nil, less or equal to 0")
	}
	if o.flushTimeout == nil || o.flushTimeout() <= 0 {
		return errors.New("flush timeout can't be nil, less or equal to 0")
	}

	return nil
}
