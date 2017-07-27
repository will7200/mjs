package http

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/will7200/mjs/apischeduler"
	"github.com/will7200/mjs/apischeduler/endpoints"
)

// NewHTTPHandler returns a handler that makes a set of endpoints available on
// predefined paths.
func NewHTTPHandler(endpoints endpoints.Endpoints) *mux.Router {
	t := mux.NewRouter()
	t.StrictSlash(true)
	m := t.PathPrefix("/api/job").Subrouter()
	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
	}
	m.Handle("/", httptransport.NewServer(
		endpoints.AddEndpoint,
		DecodeAddRequest,
		EncodeAddResponse,
		opts...,
	)).Methods("POST")
	m.Handle("/start/{id}", httptransport.NewServer(
		endpoints.StartEndpoint,
		DecodeStartRequest,
		EncodeStartResponse,
		opts...,
	)).Methods("POST")
	m.Handle("/remove/{id}", httptransport.NewServer(
		endpoints.RemoveEndpoint,
		DecodeRemoveRequest,
		EncodeRemoveResponse,
		opts...,
	)).Methods("POST")
	m.Handle("/{id}", httptransport.NewServer(
		endpoints.ChangeEndpoint,
		DecodeChangeRequest,
		EncodeChangeResponse,
		opts...,
	)).Methods("PUT")
	m.Handle("/{id}", httptransport.NewServer(
		endpoints.GetEndpoint,
		DecodeGetRequest,
		EncodeGetResponse,
		opts...,
	)).Methods("GET")
	m.Handle("/", httptransport.NewServer(
		endpoints.ListEndpoint,
		DecodeListRequest,
		EncodeListResponse,
		opts...,
	)).Methods("GET")
	m.Handle("/enable/{id}", httptransport.NewServer(
		endpoints.EnableEndpoint,
		DecodeEnableRequest,
		EncodeEnableResponse,
		opts...,
	)).Methods("POST")
	m.Handle("/disable/{id}", httptransport.NewServer(
		endpoints.DisableEndpoint,
		DecodeDisableRequest,
		EncodeDisableResponse,
		opts...,
	)).Methods("POST")
	return t
}

type errorWrapper struct {
	Error string `json:"error"`
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	msg := err.Error()
	switch err {
	case apischeduler.JobExists:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(errorWrapper{Error: msg})
}

func getHeaders(c context.Context, r *http.Request) context.Context {
	value := r.Header.Get(apischeduler.JobUniqueness)
	return context.WithValue(c, apischeduler.JobUniqueness, value)

}

// DecodeAddRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeAddRequest(_ context.Context, r *http.Request) (req interface{}, err error) {
	req = endpoints.AddRequest{}
	err = json.NewDecoder(r.Body).Decode(&r)
	return req, err
}

// EncodeAddResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeAddResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	err = e.Encode(response)
	return err
}

// DecodeStartRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeStartRequest(_ context.Context, r *http.Request) (req interface{}, err error) {
	req = endpoints.StartRequest{Id: mux.Vars(r)["id"]}
	//err = json.NewDecoder(r.Body).Decode(&r)
	return req, err
}

// EncodeStartResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeStartResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	err = e.Encode(response)
	return err
}

// DecodeRemoveRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeRemoveRequest(_ context.Context, r *http.Request) (req interface{}, err error) {
	req = endpoints.RemoveRequest{Id: mux.Vars(r)["id"]}
	//err = json.NewDecoder(r.Body).Decode(&r)
	return req, err
}

// EncodeRemoveResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeRemoveResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return err
}

// DecodeChangeRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeChangeRequest(_ context.Context, r *http.Request) (req interface{}, err error) {
	req = endpoints.ChangeRequest{Id: mux.Vars(r)["id"]}
	err = json.NewDecoder(r.Body).Decode(&r)
	return req, err
}

// EncodeChangeResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeChangeResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	err = e.Encode(response)
	return err
}

// DecodeGetRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeGetRequest(_ context.Context, r *http.Request) (req interface{}, err error) {
	req = endpoints.GetRequest{Id: mux.Vars(r)["id"]}
	//err = json.NewDecoder(r.Body).Decode(&r)
	return req, err
}

// EncodeGetResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeGetResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	err = e.Encode(response)
	return err
}

// DecodeListRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeListRequest(_ context.Context, r *http.Request) (req interface{}, err error) {
	req = endpoints.ListRequest{}
	//err = json.NewDecoder(r.Body).Decode(&r)
	return req, err
}

// EncodeListResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeListResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	err = e.Encode(response)
	return err
}

// DecodeEnableRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeEnableRequest(_ context.Context, r *http.Request) (req interface{}, err error) {
	req = endpoints.EnableRequest{Id: mux.Vars(r)["id"]}
	//err = json.NewDecoder(r.Body).Decode(&r)
	return req, err
}

// EncodeEnableResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeEnableResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	err = e.Encode(response)
	return err
}

// DecodeDisableRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeDisableRequest(_ context.Context, r *http.Request) (req interface{}, err error) {
	req = endpoints.DisableRequest{Id: mux.Vars(r)["id"]}
	//err = json.NewDecoder(r.Body).Decode(&r)
	return req, err
}

// EncodeDisableResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeDisableResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	err = e.Encode(response)
	return err
}
