package http

import (
	"net/http"

	"github.com/dee-el/go-fw/errors"
)

// Dictionary is used to store a mapping from `error.Type` to http status.
// This is to make it easy for user when err happens to its http status as business rule.
type Dictionary map[errors.Type]int

var DefaultDictionary Dictionary = Dictionary{
	errors.TypeAuthenticationError:   http.StatusUnauthorized,
	errors.TypeNotFoundError:         http.StatusNotFound,
	errors.TypeForbiddenError:        http.StatusForbidden,
	errors.TypeApplicationLimitError: http.StatusTooManyRequests,
	errors.TypeInternalServerError:   http.StatusInternalServerError,
	errors.TypeMaintenanceError:      http.StatusServiceUnavailable,
	errors.TypeBadRequestError:       http.StatusBadRequest,
}

var dictionary = DefaultDictionary

// CreateDictionary is function to create new dictionary.
// This is make user still have flexibiilty to use their own dictionary
func CreateDictionary(d Dictionary) {
	dictionary = d
}
