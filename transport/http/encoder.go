package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dee-el/go-fw/errors"
	"github.com/dee-el/go-fw/transport/http/response"
)

// JSONResponseEncoder encodes the passed response object to the HTTP response writer in JSON format.
func JSONResponseEncoder(ctx context.Context, w http.ResponseWriter, httpStatus int, res response.Response) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	return json.NewEncoder(w).Encode(res)
}

// JSONResponseEncoder encodes the passed err to client in JSON format.
// Using Dictionary to help directing err to each own HTTP status.
func JSONErrorEncoder(ctx context.Context, w http.ResponseWriter, err error) {
	// empty response, value will fill from type checking err
	// same response to standardize format response
	resp := response.NewResponse(nil, nil)

	httpStatus := http.StatusOK
	if err != nil {
		switch e := err.(type) {
		case *errors.Error:
			resp.Error = e

			httpStatus = dictionary[e.Type]
		default:
			// mask the real error
			// client(s) doesn't need to know what it is
			resp.Error = errors.ErrorInternalServer
			httpStatus = dictionary[resp.Error.Type]
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	// no need check err encoder
	json.NewEncoder(w).Encode(resp)
}
