package apischeduler

import (
	"errors"
	"fmt"
)

type contextkey string

const (
	JobUniqueness         string = "TYPE"
	APPName               string = "mjs"
	AppErrorString        string = "apischeduler.AppError"
	AppErrorFieldInternal string = "Internal"
)

var (
	JobExists error = errors.New("Job already exists")
	JobDNE    error = errors.New("Job does not exist")
	JobDBerr     error = errors.New("Interal DB Error")
)

type AppError struct {
	Internal error
	message  string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s;%s", e.Internal, e.message)
}

func GetAppError(internal error, message string) error {
	return &AppError{internal, message}
}
