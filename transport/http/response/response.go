package response

import "github.com/dee-el/go-fw/errors"

// Response is a  response sent to client in JSON format
type Response struct {
	// this field should be filled when 4xx and 5xx status returned
	Error *errors.Error `json:"error"`

	// any object from service should be on this field
	Data interface{} `json:"data"`
}

func NewResponse(data interface{}, err *errors.Error) *Response {
	return &Response{
		Data:  data,
		Error: err,
	}
}
