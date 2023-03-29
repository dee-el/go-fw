package http

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/dee-el/go-fw/tracederr"
)

func LoggedErrorHandler(logger *zap.Logger) func(r *http.Request, err error) {
	return func(r *http.Request, err error) {
		ctx := r.Context()
		// since using chi as template mux http, utilize chi context to get regexed pattern path
		rctx := chi.RouteContext(r.Context())
		path := rctx.RoutePattern()

		if path == "" {
			path = r.URL.Path
		}

		logger.Error(
			"HTTP error",
			zap.String("request_path", path),
			zap.String("request_method", r.Method),
			zap.String("request_id", chi_middleware.GetReqID(ctx)),
			zap.Any("stack_errors", tracederr.PrintErrors(err, nil)),
		)
	}
}

func NoopErrorHandler(r *http.Request, err error) {
	// simply let err gone
}
