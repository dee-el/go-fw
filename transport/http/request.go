package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/form/v4"
)

// Request is extracted request object. For querystring, user should handle it manually.
type Request struct {
	Payload interface{}
	// let it empty will be let handler retrieve it for user automatically
	URLParams URLParams
}

// this is from URL slugs, but does not rule out manual processing from RequestDecoder.
// Use can extract it via r.URL()
type URLParams map[string]string

func (up URLParams) Get(key string) string {
	val, ok := up[key]
	if ok {
		return val
	}

	return ""
}

var formDecoder = form.NewDecoder()

// RequestParser decode request to object payload.
// Only accept form JSON and form
func RequestParser(r *http.Request, payload interface{}) error {
	contents := r.Header.Get("Content-Type")
	if strings.Index(contents, "application/json") == 0 {
		return json.NewDecoder(r.Body).Decode(payload)
	} else if strings.Index(contents, "multipart/form-data") == 0 {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			return err
		}

		// ignoring file
		// any files will need manual process
		return formDecoder.Decode(payload, r.MultipartForm.Value)
	}

	err := r.ParseForm()
	if err != nil {
		return err
	}

	return formDecoder.Decode(payload, r.Form)
}
