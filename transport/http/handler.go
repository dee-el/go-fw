package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/dee-el/go-fw/errors"
	"github.com/dee-el/go-fw/transport/http/response"
)

// Note: the idea of this pattern come from `go-kit`

// RequestDecoder  extracts a user-domain request object from an HTTP request object.
type RequestDecoder func(ctx context.Context, r *http.Request) (*Request, error)

// ResponseDecoder encodes the passed response object to the HTTP response writer.
// HTTP status becomes the returned part to flexibility for user,
// such as 201 when create object, 202 to accepted the job, 401 when authentication to get nonce, etc.
type ResponseEncoder func(ctx context.Context, w http.ResponseWriter, httpStatus int, resp response.Response) error

// ErrorEncoder is responsible for encoding an error to the ResponseWriter.
// This separation response for err to easier understanding process.
type ErrorEncoder func(ctx context.Context, w http.ResponseWriter, err error)

// ErrorHandler processing err internally.
// Separation process on ErrorEncoder to make 1 purpose for each process.
type ErrorHandler func(r *http.Request, err error)

// Endpoint is the fundamental building block of servers and clients. It represents a single RPC method.
type Endpoint func(ctx context.Context, request *Request) (resp response.Response, httpStatus int, err error)

// Handler is a custom implementation of http.Handler
type Handler struct {
	endpoint        Endpoint
	requestDecoder  RequestDecoder
	responseEncoder ResponseEncoder
	errorEncoder    ErrorEncoder
	errorHandler    ErrorHandler
}

type HandlerOption func(h *Handler)

// WithResponseEncoder is an option to replace ResponseEncoder on Handler
func WithResponseEncoder(responseEncoder ResponseEncoder) HandlerOption {
	return func(h *Handler) {
		h.responseEncoder = responseEncoder
	}
}

// WithErrorEncoder is an option to replace ErrorEncoder on Handler.
func WithErrorEncoder(errorEncoder ErrorEncoder) HandlerOption {
	return func(h *Handler) {
		h.errorEncoder = errorEncoder
	}
}

// WithErrorHandler is an option to replace ErrorHandler on Handler
func WithErrorHandler(errorHandler ErrorHandler) HandlerOption {
	return func(h *Handler) {
		h.errorHandler = errorHandler
	}
}

var logger, _ = zap.NewProduction(zap.AddStacktrace(zap.PanicLevel), zap.WithCaller(false))
var logged = LoggedErrorHandler(logger)

func NewHandler(endpoint Endpoint, requestDecoder RequestDecoder, opts ...HandlerOption) *Handler {
	h := &Handler{
		endpoint:        endpoint,
		requestDecoder:  requestDecoder,
		responseEncoder: JSONResponseEncoder,
		errorEncoder:    JSONErrorEncoder,
		errorHandler:    logged,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// leverage chi context
	chiCtx := chi.RouteContext(ctx)
	urlParams := getURLParams(chiCtx)

	// input
	request, err := h.requestDecoder(ctx, r)
	if err != nil {
		_, ok := err.(*errors.Error)
		if !ok {
			// ignoring business error
			h.errorHandler(r, err)
		}

		h.errorEncoder(ctx, w, err)
		return
	}

	// overwrite when request decoder did not retrieve url params
	if len(request.URLParams) == 0 {
		request.URLParams = urlParams
	}

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			h.errorHandler(r, ctx.Err())
			h.errorEncoder(ctx, w, errors.ErrorInternalServer)
			return
		}
	default:
		// process
		response, httpStatus, err := h.endpoint(ctx, request)
		if err != nil {
			h.errorHandler(r, err)
			h.errorEncoder(ctx, w, err)
			return
		}

		// output
		err = h.responseEncoder(ctx, w, httpStatus, response)
		if err != nil {
			h.errorHandler(r, err)
			h.errorEncoder(ctx, w, err)
			return
		}
	}
}

func getURLParams(chiCtx *chi.Context) URLParams {
	params := URLParams{}
	for _, key := range chiCtx.URLParams.Keys {
		params[key] = chiCtx.URLParam(key)
	}

	return params
}
