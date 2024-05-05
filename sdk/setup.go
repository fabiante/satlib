package sdk

type setupOption func(*setup) error

type setup struct {
}

func Setup(opts ...setupOption) *setup {
	s := &setup{}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			panic(err)
		}
	}

	return s
}

func (s *setup) Run(handler func(ctx *Context) error) {
	run(handler)
}
