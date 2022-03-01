package idgenerator

type Option func(cfg *config)

type config struct {
	// instanceId is the host id.
	instanceId *int64
}

func applyOptions(opts ...Option) config {
	cfg := config{
		instanceId: nil,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.instanceId == nil {
		x := defaultInstanceID()
		cfg.instanceId = &x
	}
	return cfg
}

// WithInstanceId specified caller's instance id.
func WithInstanceId(n int64) Option {
	return func(cfg *config) {
		cfg.instanceId = &n
	}
}
