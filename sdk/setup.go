package sdk

import (
	"fmt"
	"log/slog"
	"net/http"
)

type setup struct {
	httpHandlers map[string]http.HandlerFunc
}

func Setup(opts ...SetupOption) *setup {
	s := &setup{
		httpHandlers: make(map[string]http.HandlerFunc),
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			panic(err)
		}
	}

	return s
}

func (s *setup) Run(handler func(ctx *HttpContext) error) {
	addr := fmt.Sprintf(":%v", getEnv("PORT", "8080"))

	slog.Info("SDK listening", "addr", addr)

	mux := http.NewServeMux()

	// Add main handler to handler set
	s.httpHandlers["/"] = func(writer http.ResponseWriter, request *http.Request) {
		ctx := newContext(request, writer)
		if err := handler(ctx); err != nil {
			panic(err)
		}
	}

	// Register http handlers
	for path, handler := range s.httpHandlers {
		mux.HandleFunc(path, handler)
	}

	err := http.ListenAndServe(addr, mux)

	if err != nil {
		panic(err)
	}
}
