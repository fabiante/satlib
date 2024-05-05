package sdk

import (
	"fmt"
	"log/slog"
	"net/http"
)

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
	addr := fmt.Sprintf(":%v", getEnv("PORT", "8080"))

	slog.Info("SDK listening", "addr", addr)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		ctx := newContext(request, writer)
		if err := handler(ctx); err != nil {
			panic(err)
		}
	})

	err := http.ListenAndServe(addr, mux)

	if err != nil {
		panic(err)
	}
}
