package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	statusOk    = "OK"
	statusError = "Error"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func OK() Response {
	return Response{
		Status: statusOk,
	}
}

func Error(msg string) Response {
	return Response{
		Status: statusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		errMsgs = append(errMsgs, fmt.Sprintf(`field %s is a required field`, err.Field()))
	}

	return Response{
		Status: statusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
