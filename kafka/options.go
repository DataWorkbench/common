package kafka

import (
	"time"
)

type Option func(o *Options)

type Options struct {
	// option for consumerHandler.
	batchMode     bool
	batchMax      int
	retryInterval time.Duration
}

func applyOptions(options ...Option) Options {
	opts := Options{
		batchMode:     false,
		batchMax:      256,
		retryInterval: time.Second * 5,
	}

	for _, option := range options {
		option(&opts)
	}
	return opts
}

// WithBatchMode controls the consumerHandler whether enable the `batchMode`
//
// If true, the behavior of consumerHandler will try to consume as many messages at once.
// Otherwise, the behavior of consumerHandler will consume one message at once.
func WithBatchMode(ok bool) Option {
	return func(o *Options) {
		o.batchMode = ok
	}
}

// WithBatchMax sets the maximum messages of consumed at once if `batchMode` is enabled.
// Defaults 128.
func WithBatchMax(max int) Option {
	return func(o *Options) {
		o.batchMax = max
	}
}

// RetryInterval sets the retry interval time when consumerHandler returns error.
// Defaults 5s.
func WithRetryInterval(d time.Duration) Option {
	return func(o *Options) {
		o.retryInterval = d
	}
}
