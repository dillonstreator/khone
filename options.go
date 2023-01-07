package khone

type config struct {
	concurrency int
}

func newConfig(options ...option) *config {
	cfg := &config{
		concurrency: 1,
	}

	for _, o := range options {
		o(cfg)
	}

	return cfg
}

type option func(cfg *config)

func WithConcurrency(concurrency int) option {
	return func(cfg *config) {
		cfg.concurrency = concurrency
	}
}
