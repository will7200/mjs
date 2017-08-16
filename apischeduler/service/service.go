package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/will7200/mjs/apischeduler"
	"github.com/will7200/mjs/job"
)

// Implement yor service methods methods.
// e.x: Foo(ctx context.Context,s string)(rs string, err error)
type APISchedulerService interface {
	//METHODS: POST
	//PATH: /
	Add(ctx context.Context, reqjob job.Job) (id string, err error)
	//METHODS: POST
	//PATH: /start/{id}
	Start(ctx context.Context, id string) (message string, err error)
	//METHODS: POST
	//PATH: /remove/{id}
	Remove(ctx context.Context, id string) (message string, err error)
	//METHODS: PUT
	//PATH: /change/{id}
	Change(ctx context.Context, id string, reqjob job.Job) (message string, err error)
	//METHODS: GET
	//PATH: /{id}
	Get(ctx context.Context, id string) (result *job.Job, err error)
	//METHODS: GET
	//PATH: /
	List(ctx context.Context) (results *[]job.Job, err error)
	//METHODS: POST
	//PATH: /enable
	Enable(ctx context.Context, id string) (message string, err error)
	//METHODS: POST
	//PATH: /disable
	Disable(ctx context.Context, id string) (message string, err error)
}

type stubAPISchedulerService struct {
	db       *gorm.DB
	dispatch *job.Dispatcher
}

// Get a new instance of the service.
// If you want to add service middleware this is the place to put them.
func New(db *gorm.DB, dispatch *job.Dispatcher) (s APISchedulerService) {
	s = &stubAPISchedulerService{db, dispatch}
	return s
}

// Implement the business logic of Add
func (ap *stubAPISchedulerService) Add(ctx context.Context, reqjob job.Job) (id string, err error) {
	db := ap.db
	j := &reqjob
	if val, ok := ctx.Value(apischeduler.JobUniqueness).(string); val != "" && ok {
		if val == "UNIQUE" {
			jj := &job.Job{}
			if !db.Where(job.Job{Domain: j.Domain, SubDomain: j.SubDomain, Name: j.Name,
				Application: j.Application}).First(jj).RecordNotFound() {
				err = apischeduler.GetAppError(apischeduler.JobExists, jj.ID)
				return
			}
		}
	}
	if err = j.ParseSchedule(); err != nil {
		return
	}
	if err = db.Create(j).Error; err != nil {
		return
	}
	j.StartWaiting(ap.dispatch)
	id = j.ID
	return id, err
}

// Implement the business logic of Start
func (ap *stubAPISchedulerService) Start(ctx context.Context, id string) (message string, err error) {
	d, err := ap.Get(ctx, id)
	if err != nil {
		return "", err
	}
	ap.dispatch.AddFutureJob(d, time.Millisecond*500)
	if err != nil {
		return "", err
	}
	message = fmt.Sprintf("Job with id %s has been started", d.ID)
	return message, err
}

// Implement the business logic of Remove
func (ap *stubAPISchedulerService) Remove(ctx context.Context, id string) (message string, err error) {
	d := &job.Job{}
	d, err = ap.Get(ctx, id)
	if err != nil {
		return "", err
	}
	if err = ap.db.Delete(d).Error; err != nil {
		err = apischeduler.GetAppError(apischeduler.JobDBerr, err.Error())
		return "", err
	}
	message = fmt.Sprintf("Job with id %s has been removed", d.ID)
	return message, err
}

// Implement the business logic of Change
func (ap *stubAPISchedulerService) Change(ctx context.Context, id string, reqjob job.Job) (message string, err error) {
	d, err := ap.Get(ctx, id)
	if err != nil {
		return "", err
	}
	if err = ap.db.Model(d).Update(reqjob).Error; err != nil {
		err = apischeduler.GetAppError(fmt.Errorf("Cannot Update record with id %s", id), err.Error())
		return "", err
	}
	ap.dispatch.RemoveWorkRequest(d)
	d.ParseSchedule()
	d.StartWaiting(ap.dispatch)
	//TODO : MAYBE IMPLEMENT TO GET THE AMOUNT OF FIELDS CHANGED
	message = fmt.Sprintf("Job with id %s has changed", d.ID)
	return message, err
}

// Implement the business logic of Get
func (ap *stubAPISchedulerService) Get(ctx context.Context, id string) (result *job.Job, err error) {
	r := &job.Job{}
	if err = ap.db.Where(job.Job{ID: id}).First(r).Error; err != nil {
		err = fmt.Errorf("Error: %s;\nDatabaseError:%s", apischeduler.JobDNE, err.Error())
		return nil, err
	}
	result = r
	return result, err
}

// Implement the business logic of List
func (ap *stubAPISchedulerService) List(ctx context.Context) (results *[]job.Job, err error) {
	r := &[]job.Job{}
	if err = ap.db.Find(r).Error; err != nil {
		return
	}
	results = r
	return results, err
}

// Implement the business logic of Enable
func (ap *stubAPISchedulerService) Enable(ctx context.Context, id string) (message string, err error) {
	d, err := ap.Get(ctx, id)
	if err != nil {
		return "", err
	}
	if err := ap.db.Model(d).Update(job.Job{IsActive: true}).Error; err != nil {
		err = apischeduler.GetAppError(fmt.Errorf("Cannot Enable record with id %s", id), err.Error())
		return "", err
	}
	message = fmt.Sprintf("Job with id %s has been enabled", d.ID)
	return message, err
}

// Implement the business logic of Disable
func (ap *stubAPISchedulerService) Disable(ctx context.Context, id string) (message string, err error) {
	d, err := ap.Get(ctx, id)
	if err != nil {
		return "", err
	}
	if err = ap.db.Model(d).Update(job.Job{IsActive: false}).Error; err != nil {
		err = apischeduler.GetAppError(fmt.Errorf("Cannot Disable record with id %s", id), err.Error())
		return "", err
	}
	message = fmt.Sprintf("Job with id %s has been disabled", d.ID)
	return message, err
}
