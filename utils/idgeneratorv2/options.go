package idgeneratorv2

type Option func(cfg *config)

type config struct {
	// instanceId is the host id.
	instanceId *int64

	// hashSalt is the secret used to make the generated id harder to guess
	hashSalt string

	// hashMinLength is the minimum length of a generated id
	hashMinLength int

	// hashAlphabet is the alphabet used to generate new ids
	hashAlphabet string
}

func applyOptions(opts ...Option) config {
	cfg := config{
		instanceId:    nil,
		hashSalt:      "dataomnis",
		hashMinLength: 16,
		hashAlphabet:  "abcdefghijklmnopqrstuvwxyz1234567890",
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

// WithHashSalt set the salt.
func WithHashSalt(salt string) Option {
	return func(cfg *config) {
		cfg.hashSalt = salt
	}
}

// WithHashMinLength set the hashMinLength.
func WithHashMinLength(n int) Option {
	return func(cfg *config) {
		cfg.hashMinLength = n
	}
}

// WithHashAlphabet set the hashAlphabet.
func WithHashAlphabet(s string) Option {
	return func(cfg *config) {
		cfg.hashAlphabet = s
	}
}
