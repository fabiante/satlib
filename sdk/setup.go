package sdk

import (
	"fmt"
	"log/slog"
	"net/http"
)

type SetupOption func(*setup) error

type setup struct {
	handlers map[string]http.HandlerFunc
}

func Setup(opts ...SetupOption) *setup {
	s := &setup{
		handlers: make(map[string]http.HandlerFunc),
	}

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

	// Register custom handlers
	for path, handler := range s.handlers {
		mux.HandleFunc(path, handler)
	}

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

func withHandler(path string, handler http.HandlerFunc) SetupOption {
	return func(s *setup) error {
		s.handlers[path] = handler
		return nil
	}
}

func WithOpenApiSpec(spec []byte) SetupOption {
	return withHandler("/openapi.yml", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/yaml")
		_, _ = writer.Write(spec)
	})
}
