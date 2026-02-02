package errs

import (
	"net/http"
)

// HTTP Status Code: 401 ("UNAUTHORIZED")
func NewUnauthorizedError(message string, overide bool) *HTTPError {
	return &HTTPError{
		Code:    MakeUpperCaseWithUnderscores(http.StatusText(http.StatusUnauthorized)),
		Message: message,
		Status:  http.StatusUnauthorized,
		Overide: overide,
	}
}

// HTTP Status Code: 403 ("FORBIDDEN")
func NewForbiddenError(message string, overide bool) *HTTPError {
	return &HTTPError{
		Code:    MakeUpperCaseWithUnderscores(http.StatusText(http.StatusForbidden)),
		Message: message,
		Status:  http.StatusForbidden,
		Overide: overide,
	}
}

// HTTP Status Code: 400 ("BAD REQUEST")
func NewBadRequestError(message string, overide bool, action *Action, errors []FieldError) *HTTPError {
	return &HTTPError{
		Code:    http.StatusText(http.StatusBadRequest),
		Message: message,
		Status:  http.StatusBadRequest,
		Overide: overide,
		Errors:  errors,
		Action:  action,
	}
}

// HTTP Status Code: 404 ("NOT FOUND")
func NewNotFoundError(message string, overide bool, code *string) *HTTPError {
	formattedCode := MakeUpperCaseWithUnderscores(http.StatusText(http.StatusNotFound))

	if code != nil {
		formattedCode = *code
	}

	return &HTTPError{
		Code:    formattedCode,
		Message: message,
		Status:  http.StatusNotFound,
		Overide: overide,
	}
}

// HTTP Status Code: 500 ("INTERNAL SERVER ERROR")
func NewInternalServerError() *HTTPError {
	return &HTTPError{
		Code:    MakeUpperCaseWithUnderscores(http.StatusText(http.StatusInternalServerError)),
		Message: http.StatusText(http.StatusInternalServerError),
		Status:  http.StatusInternalServerError,
		Overide: false,
	}
}

func ValidationError(err error) *HTTPError {
	return NewBadRequestError("Validation failed: "+err.Error(), false, nil, nil)
}
