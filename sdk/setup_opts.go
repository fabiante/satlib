package sdk

import "net/http"

type SetupOption func(*setup) error

func withHandler(path string, handler http.HandlerFunc) SetupOption {
	return func(s *setup) error {
		s.httpHandlers[path] = handler
		return nil
	}
}

func WithOpenApiSpec(spec []byte) SetupOption {
	return withHandler("/openapi.yml", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/yaml")
		_, _ = writer.Write(spec)
	})
}
