package idgenerator

type Option func(cfg *config)

type config struct {
	instanceId *int64
}

func applyOptions(opts ...Option) config {
	cfg := config{}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

// WithInstanceId specified caller's instance id.
func WithInstanceId(n int64) Option {
	return func(cfg *config) {
		cfg.instanceId = &n
	}
}
