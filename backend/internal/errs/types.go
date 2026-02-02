package errs

import (
	"strings"
)

/**
@dev Global Error Structs used over application
*/

// Deals with form based errors
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// this action related info is used to inform client to take action
// and, also to redirect user to other specific page
type ActionType string

const (
	ActionTypeRedirect ActionType = "redirect"
)

type Action struct {
	Type    ActionType `json:"type"`
	Message string     `json:"message"`
	Value   string     `json:"value"`
}

// this is final error struct we sent to client
type HTTPError struct {
	Code    string `json:"code"` // Ex: TODO_NOT_FOUND type of error code
	Message string `json:"message"`
	Status  int    `json:"status"` // HTTP status code
	Overide bool   `json:"override"`
	// field/form level errors
	Errors []FieldError `json:"errors"`
	// action to be taken
	Action *Action `json:"action"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) Is(target error) bool {
	_, ok := target.(*HTTPError)

	return ok
}

func (e *HTTPError) WithMessage(message string) *HTTPError {
	return &HTTPError{
		Code:    e.Code,
		Message: message,
		Status:  e.Status,
		Overide: e.Overide,
		Errors:  e.Errors,
		Action:  e.Action,
	}
}

func MakeUpperCaseWithUnderscores(str string) string {
	return strings.ToUpper(strings.ReplaceAll(str, " ", "_"))
}
