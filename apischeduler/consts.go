package apischeduler

import "errors"

type contextkey string

const (
	JobUniqueness string = "TYPE"
)

var (
	JobExists error = errors.New("Job already exists")
	JobDNE    error = errors.New("Job does not exist")
)
