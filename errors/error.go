package errors

type Code int

type Type string

// Error is a standard when something wrong happen on service.
// This Error happened as business error, any other errors should be wrapped by ErrorWithStack or let it reach to user.
//
// Idea for this structure is some error will drop on same type.
// example:
// - TypeAuthenticationError will have code 101, 105, 106, 109, 110
// example representation:
//   - 101: Invalid access token
//   - 105: Username of password is wrong
//   - 106: Access token is expired
//   - 109: Refresh token is not found
//   - 110: API key is wrong
// or whatever user want, this thing will never be exhausted because user can craete their own combination Type and Code.
// Another reason is easier to tracking flow of business.
type Error struct {
	Type    Type   `json:"type"`
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Fields  Fields `json:"fields"`
}

type Fields map[string]string

func (e *Error) AddField(k, v string) {
	if len(e.Fields) == 0 {
		e.Fields = map[string]string{}
	}

	e.Fields[k] = v
}

// Error satisfying interface Error, just return the message
func (e *Error) Error() string {
	return e.Message
}

func New(t Type, c Code, msg string) *Error {
	return &Error{
		Code:    c,
		Type:    t,
		Message: msg,
		Fields:  nil,
	}
}

// Reserved codes, so user does not set this again
const (
	// 1xx
	Code101 Code = 101

	// 8xx
	Code805 = 805
	Code822 = 822
	Code825 = 825
	Code831 = 831

	// 9xx
	Code901 = 901
	Code910 = 910
)

// Reserved types
const (
	TypeAuthenticationError   Type = "AuthenticationError"
	TypeNotFoundError         Type = "NotFoundError"
	TypeInternalServerError   Type = "InternalServerError"
	TypeForbiddenError        Type = "ForbiddenError"        // permission denied
	TypeApplicationLimitError Type = "ApplicationLimitError" // throttle
	TypeMaintenanceError      Type = "MaintenanceError"
	TypeBadRequestError       Type = "BadRequestError"
)

// Reserved errors
// any other business should be added on its own error dictionary
var (
	ErrorAuthentication   = New(TypeAuthenticationError, Code805, "Authentication failed")
	ErrorNotFound         = New(TypeNotFoundError, Code822, "Resource not found")
	ErrorForbidden        = New(TypeForbiddenError, Code825, "Permission denied")
	ErrorInternalServer   = New(TypeInternalServerError, Code901, "Oops, something went wrong")
	ErrorMaintenance      = New(TypeMaintenanceError, Code910, "Sorry, app is under maintenance")
	ErrorApplicationLimit = New(TypeApplicationLimitError, Code831, "Application limit is exceeded")
	ErrorBadRequest       = New(TypeBadRequestError, Code101, "Bad request")
)
