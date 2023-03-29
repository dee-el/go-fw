package middleware

import (
	"context"
	"net/http"
)

type ctxKey string

const apiVersionCtxKey ctxKey = "api.version"

func APIVersion(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), apiVersionCtxKey, version))
			next.ServeHTTP(w, r)
		})
	}
}
