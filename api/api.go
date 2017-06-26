package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/will7200/mjs/job"
)

var (
	APIUrlPrefix = "/api/"

	JobPath    = "job/"
	APIJobPath = APIUrlPrefix + JobPath

	contentType     = "Content-Type"
	jsonContentType = "application/json;charset=UTF-8"
)

var d *job.Dispatcher
var db *gorm.DB

func unmarshalNewJob(r *http.Request) (*job.Job, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Errorf("Error occured when reading r.Body: %s", err)
		return nil, err
	}
	defer r.Body.Close()

	return job.RawJob(body, d)
}

type addJobResponse struct {
	ID string `json:"id"`
}

func HandleAddJob(w http.ResponseWriter, r *http.Request) {
	j, err := unmarshalNewJob(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if value := r.Header.Get("TYPE"); value != "" {
		if value == "UNIQUE" {
			jj := &job.Job{}
			if !db.Where(job.Job{Domain: j.Domain, SubDomain: j.SubDomain, Name: j.Name,
				Application: j.Application}).First(jj).RecordNotFound() {
				http.Error(w, errors.New("Job already exists").Error(), http.StatusBadRequest)
				return
			}
		}
	}
	if db.Create(j).Error != nil {
		http.Error(w, errors.New("Update to Save Job").Error(), http.StatusBadRequest)
		return
	}
	resp := &addJobResponse{
		ID: j.ID,
	}

	w.Header().Set(contentType, jsonContentType)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Errorf("Error occured when marshalling response: %s", err)
		return
	}
}
func HandleJob(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	j := &job.Job{}
	if err := db.Where(job.Job{ID: id}).First(j).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Could not find Job with id %s", err.Error())
		return
	}
	if j == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if r.Method == "DELETE" {
		err := db.Delete(j).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	} else if r.Method == "GET" {
		handleGetJob(w, r, j)
	}
}

type JobResponse struct {
	Job *job.Job `json:"job"`
}
type JobsResponse struct {
	Job []job.Job `json:"jobs"`
}

func handleGetJob(w http.ResponseWriter, r *http.Request, j *job.Job) {
	resp := &JobResponse{
		Job: j,
	}

	encodeJSONResponse(w, r, resp)
}

func HandleListJobs(w http.ResponseWriter, r *http.Request) {
	resp := &JobsResponse{}
	j := []job.Job{}
	if err := db.Find(&j).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Interal Error %s", err.Error())
		return
	}
	resp.Job = j
	encodeJSONResponse(w, r, resp)
}
func HandleStartJob(w http.ResponseWriter, r *http.Request) {

}
func HandleEnableJob(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	j := &job.Job{}
	if err := db.Where(job.Job{ID: id}).First(j).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Could not find Job with id %s", err.Error())
		return
	}
	if j == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	j.IsActive = true
	db.Save(j)
}
func HandleDisableJob(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	j := &job.Job{}
	if err := db.Where(job.Job{ID: id}).First(j).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Could not find Job with id %s", err.Error())
		return
	}
	if j == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	j.IsActive = false
	db.Save(j)
}
func HandleListJobStats(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	j := &job.Job{}
	if err := db.Where(job.Job{ID: id}).First(j).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Could not find Job with id %s", err.Error())
		return
	}
	if j == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	s := &[]job.JobStats{}
	db.Model(j).Association("Stats").Find(s)
	encodeJSONResponse(w, r, s)
}
func SetupAPIRoutes(r *mux.Router) {
	// Route for creating a job
	r.HandleFunc(APIJobPath, HandleAddJob).Methods("POST")
	// Route for deleting and getting a job
	r.HandleFunc(APIJobPath+"{id}/", HandleJob).Methods("DELETE", "GET")
	// Route for getting job stats
	r.HandleFunc(APIJobPath+"stats/{id}/", HandleListJobStats).Methods("GET")
	// Route for listing all jobs
	r.HandleFunc(APIJobPath, HandleListJobs).Methods("GET")
	// Route for manually start a job
	r.HandleFunc(APIJobPath+"start/{id}/", HandleStartJob).Methods("POST")
	// Route for manually enable a job
	r.HandleFunc(APIJobPath+"enable/{id}/", HandleEnableJob).Methods("POST")
	// Route for manually disable a job
	r.HandleFunc(APIJobPath+"disable/{id}/", HandleDisableJob).Methods("POST")
}
func StartServer(listenAddr string, dd *job.Dispatcher, ddb *gorm.DB) error {
	r := mux.NewRouter()
	d = dd
	db = ddb
	// Allows for the use for /job as well as /job/
	r.StrictSlash(true)
	SetupAPIRoutes(r)
	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 7 * time.Second,
		Addr:         listenAddr,
		Handler:      r,
	}
	return server.ListenAndServe()
}

func encodeJSONResponse(w http.ResponseWriter, r *http.Request, j interface{}) {
	w.Header().Set(contentType, jsonContentType)
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	if err := e.Encode(j); err != nil {
		log.Errorf("Error occured when marshalling response: %s", err)
		return
	}
}
