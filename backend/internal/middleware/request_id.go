package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// name of the header storing request ID and key of the request ID
const (
	RequestIDHeader = "X-Request-ID"
	RequestIDKey    = "request_id"
)

// a function to process middleware. This is our middleware function
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// using context we get access to different info. about our request(headers,body,query params, payloads, etc)
		return func(c echo.Context) error {
			requestID := c.Request().Header.Get(RequestIDHeader)
			if requestID == "" {
				requestID = uuid.New().String() // 4c90fc3f-39cc-4b04-af21-c83ee64aa67e
			}

			// saves data in the context
			c.Set(RequestIDKey, requestID)
			c.Response().Header().Set(RequestIDHeader, requestID)

			// this function hands over the execution to next middleware
			return next(c)
		}
	}
}

// getter function to get request ID from context
func GetRequestID(c echo.Context) string {
	if requestID, ok := c.Get(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}
