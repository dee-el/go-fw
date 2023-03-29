package http

import (
	"net/http"
	"strings"
	"time"

	chi "github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	chi_cors "github.com/go-chi/cors"

	"github.com/dee-el/go-fw/transport/http/middleware"
)

// Server work as root http.Handler.
// Use `chi` as mux, with that mind format URL slugs will be strictly following `chi` pattern.
// Example: `/users/{id}`
//
// Most handler registration functions using type `Handler`, a custom http.Handler.
// Within `Handler`, the response will always return in JSON format.
type Server struct {
	mux                   *chi.Mux
	enableBasicMiddleware bool
	timeoutInSecond       time.Duration
	// metrics               *middleware.Metrics
	tracing *middleware.Tracing
}

// NewServer returns new Server instance
func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		timeoutInSecond:       time.Second * time.Duration(60),
		enableBasicMiddleware: true,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.mux = s.init()
	return s
}

type ServerOption func(*Server)

func WithTimeoutInSecond(tm time.Duration) ServerOption {
	return func(s *Server) {
		s.timeoutInSecond = time.Second * tm
	}
}

func WithToggleBasicMiddleware(t bool) ServerOption {
	return func(s *Server) {
		s.enableBasicMiddleware = t
	}
}

// func WithMetrics(m *middleware.Metrics) ServerOption {
// 	return func(s *Server) {
// 		s.metrics = m
// 	}
// }

func WithTracing(t *middleware.Tracing) ServerOption {
	return func(s *Server) {
		s.tracing = t
	}
}

// Server returns a http.Handler.
func (s *Server) Handler() http.Handler {
	return s.mux
}

// Middleware is function to add other middleware
func (s *Server) Middleware(fn func(next http.Handler) http.Handler) {
	s.mux.Use(fn)
}

func (s *Server) Get(path string, hn *Handler) {
	s.mux.Get(path, hn.ServeHTTP)
}

func (s *Server) Head(path string, hn *Handler) {
	s.mux.Head(path, hn.ServeHTTP)
}

func (s *Server) Post(path string, hn *Handler) {
	s.mux.Post(path, hn.ServeHTTP)
}

func (s *Server) Put(path string, hn *Handler) {
	s.mux.Put(path, hn.ServeHTTP)
}

func (s *Server) Patch(path string, hn *Handler) {
	s.mux.Patch(path, hn.ServeHTTP)
}

func (s *Server) Delete(path string, hn *Handler) {
	s.mux.Delete(path, hn.ServeHTTP)
}

func (s *Server) Connect(path string, hn *Handler) {
	s.mux.Connect(path, hn.ServeHTTP)
}

func (s *Server) Options(path string, hn *Handler) {
	s.mux.Options(path, hn.ServeHTTP)
}

func (s *Server) Method(method, path string, hn *Handler) {
	s.mux.Method(strings.ToUpper(method), path, hn)
}

func (s *Server) Mount(path string, sub *Server) {
	s.mux.Mount(path, sub.Handler())
}

// MethodFunc is custom handler registration function.
// User should use this if they want return other format instead of JSON.
func (s *Server) MethodFunc(method, path string, hn http.Handler) {
	s.mux.Method(strings.ToUpper(method), path, hn)
}

func (s *Server) Route(path string, fn func(sub *Server)) {
	// every basic middleware should follow root / parent handler, no need initiate it multiple times on every child handler(s)
	// well, user can do it tho if they want
	sub := NewServer(WithToggleBasicMiddleware(false))
	fn(sub)

	s.Mount(path, sub)
}

// all this middleware should be set even if some of these middleware are not needed
func (s *Server) init() *chi.Mux {
	mux := chi.NewRouter()

	if s.enableBasicMiddleware {
		mux = basicMiddleware(mux)

		// Set a timeout value on the request context (ctx), that will signal
		// through ctx.Done() that the request has timed out and further
		// processing should be stopped.
		mux.Use(chi_middleware.Timeout(s.timeoutInSecond))
	}

	// if s.metrics != nil {
	// 	mux.Use(s.metrics.Handler)
	// }

	if s.tracing != nil {
		mux.Use(s.tracing.Handler)
	}

	return mux
}

func basicMiddleware(mux *chi.Mux) *chi.Mux {
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := chi_cors.New(chi_cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           120, // Maximum value not ignored by any of major browsers
	})
	mux.Use(cors.Handler)

	mux.Use(chi_middleware.RequestID)
	mux.Use(chi_middleware.RealIP)

	mux.Use(middleware.Recoverer)

	return mux
}
