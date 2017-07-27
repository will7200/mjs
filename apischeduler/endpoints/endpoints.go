package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/will7200/mjs/apischeduler/service"
	"github.com/will7200/mjs/job"
)

// Endpoints collects all of the endpoints that compose an add service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.

type Endpoints struct {
	AddEndpoint     endpoint.Endpoint
	StartEndpoint   endpoint.Endpoint
	RemoveEndpoint  endpoint.Endpoint
	ChangeEndpoint  endpoint.Endpoint
	GetEndpoint     endpoint.Endpoint
	ListEndpoint    endpoint.Endpoint
	EnableEndpoint  endpoint.Endpoint
	DisableEndpoint endpoint.Endpoint
}
type AddRequest struct {
	Reqjob job.Job
}
type AddResponse struct {
	Id  string
	Err error `json:",omitempty"`
}
type StartRequest struct {
	Id string
}
type StartResponse struct {
	Message string
	Err     error `json:",omitempty"`
}
type RemoveRequest struct {
	Id string
}
type RemoveResponse struct {
	Message string
	Err     error `json:",omitempty"`
}
type ChangeRequest struct {
	Id     string
	Reqjob job.Job
}
type ChangeResponse struct {
	Message string
	Err     error `json:",omitempty"`
}
type GetRequest struct {
	Id string
}
type GetResponse struct {
	Result *job.Job
	Err    error `json:",omitempty"`
}
type ListRequest struct{}
type ListResponse struct {
	Results *[]job.Job
	Err     error `json:",omitempty"`
}
type EnableRequest struct {
	Id string
}
type EnableResponse struct {
	Message string
	Err     error `json:",omitempty"`
}
type DisableRequest struct {
	Id string
}
type DisableResponse struct {
	Message string
	Err     error `json:",omitempty"`
}

func New(svc service.APISchedulerService) (ep Endpoints) {
	ep.AddEndpoint = MakeAddEndpoint(svc)
	ep.StartEndpoint = MakeStartEndpoint(svc)
	ep.RemoveEndpoint = MakeRemoveEndpoint(svc)
	ep.ChangeEndpoint = MakeChangeEndpoint(svc)
	ep.GetEndpoint = MakeGetEndpoint(svc)
	ep.ListEndpoint = MakeListEndpoint(svc)
	ep.EnableEndpoint = MakeEnableEndpoint(svc)
	ep.DisableEndpoint = MakeDisableEndpoint(svc)
	return ep
}

// MakeAddEndpoint returns an endpoint that invokes Add on the service.
// Primarily useful in a server.
func MakeAddEndpoint(svc service.APISchedulerService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRequest)
		id, err := svc.Add(ctx, req.Reqjob)
		return AddResponse{Id: id, Err: err}, err
	}
}

// MakeStartEndpoint returns an endpoint that invokes Start on the service.
// Primarily useful in a server.
func MakeStartEndpoint(svc service.APISchedulerService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StartRequest)
		message, err := svc.Start(ctx, req.Id)
		return StartResponse{Message: message, Err: err}, err
	}
}

// MakeRemoveEndpoint returns an endpoint that invokes Remove on the service.
// Primarily useful in a server.
func MakeRemoveEndpoint(svc service.APISchedulerService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RemoveRequest)
		message, err := svc.Remove(ctx, req.Id)
		return RemoveResponse{Message: message, Err: err}, err
	}
}

// MakeChangeEndpoint returns an endpoint that invokes Change on the service.
// Primarily useful in a server.
func MakeChangeEndpoint(svc service.APISchedulerService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangeRequest)
		message, err := svc.Change(ctx, req.Id, req.Reqjob)
		return ChangeResponse{Message: message, Err: err}, err
	}
}

// MakeGetEndpoint returns an endpoint that invokes Get on the service.
// Primarily useful in a server.
func MakeGetEndpoint(svc service.APISchedulerService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)
		result, err := svc.Get(ctx, req.Id)
		return GetResponse{Result: result, Err: err}, err
	}
}

// MakeListEndpoint returns an endpoint that invokes List on the service.
// Primarily useful in a server.
func MakeListEndpoint(svc service.APISchedulerService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		results, err := svc.List(ctx)
		return ListResponse{Results: results, Err: err}, err
	}
}

// MakeEnableEndpoint returns an endpoint that invokes Enable on the service.
// Primarily useful in a server.
func MakeEnableEndpoint(svc service.APISchedulerService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(EnableRequest)
		message, err := svc.Enable(ctx, req.Id)
		return EnableResponse{Message: message, Err: err}, err
	}
}

// MakeDisableEndpoint returns an endpoint that invokes Disable on the service.
// Primarily useful in a server.
func MakeDisableEndpoint(svc service.APISchedulerService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DisableRequest)
		message, err := svc.Disable(ctx, req.Id)
		return DisableResponse{Message: message, Err: err}, err
	}
}
