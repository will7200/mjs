package apischeduler

import "errors"

type contextkey string

const (
	JobUniqueness string = "TYPE"
	APPName       string = "mjs"
)

var (
	JobExists error = errors.New("Job already exists")
	JobDNE    error = errors.New("Job does not exist")
)
