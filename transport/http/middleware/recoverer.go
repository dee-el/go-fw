package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/dee-el/go-fw/errors"
	"github.com/dee-el/go-fw/transport/http/response"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				// logging what make panic happens
				//log.New().Panic(rvr, debug.Stack())

				// rewrite the response
				resp := response.NewResponse(nil, errors.ErrorInternalServer)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				// no need check err encoder
				json.NewEncoder(w).Encode(resp)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
