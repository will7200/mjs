package service

import (
	log "github.com/sirupsen/logrus"
	"context"
	"github.com/will7200/mjs/job"
	"time"
	"fmt"
)

type Middleware func(APISchedulerService) APISchedulerService
type loggingMiddleware struct {
	logger *log.Logger
	next   APISchedulerService
}

func (mw loggingMiddleware) Add(ctx context.Context, reqjob job.Job) (id string, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(log.Fields{
			"id": id,
			"err": fmt.Sprintf("%v", err),
			"took": time.Since(begin),
		}).Info("Starting job api")
	}(time.Now())

	id, err = mw.next.Add(ctx, reqjob)
	return
}

func (mw loggingMiddleware) Start(ctx context.Context, id string) (message string, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(log.Fields{
			"id": id,
			"message": message,
			"err": err,
			"took": time.Since(begin),
		}).Info("Starting job api")
	}(time.Now())

	message, err = mw.next.Start(ctx, id)
	return
}

func (mw loggingMiddleware) Remove(ctx context.Context, id string) (message string, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(log.Fields{
			"id": id,
			"message": message,
			"err": fmt.Sprintf("%v", err),
			"took": time.Since(begin),
		}).Info("Removing job api")
	}(time.Now())

	message, err = mw.next.Remove(ctx, id)
	return
}

func (mw loggingMiddleware) Change(ctx context.Context, id string, reqjob job.Job) (message string, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(log.Fields{
			"id": id,
			"message": message,
			"err": fmt.Sprintf("%v", err),
			"took": time.Since(begin),
		}).Info("Changing job api")
	}(time.Now())

	message, err = mw.next.Change(ctx, id, reqjob)
	return
}

func (mw loggingMiddleware) Get(ctx context.Context, id string) (result *job.Job, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(log.Fields{
			"id": id,
			"err": fmt.Sprintf("%v", err),
			"took": time.Since(begin),
		}).Info("Get job api")
	}(time.Now())

	result, err = mw.next.Get(ctx, id)
	return
}

func (mw loggingMiddleware) List(ctx context.Context) (results *[]job.Job, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(log.Fields{
			"err": err,
			"took": time.Since(begin),
		}).Info("Listing jobs api")
	}(time.Now())

	results, err = mw.next.List(ctx)
	return
}

func (mw loggingMiddleware) Enable(ctx context.Context, id string) (message string, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(log.Fields{
			"id": id,
			"err": fmt.Sprintf("%v", err),
			"took": time.Since(begin),
		}).Info("Enable job api")
	}(time.Now())

	message, err = mw.next.Enable(ctx, id)
	return
}

func (mw loggingMiddleware) Disable(ctx context.Context, id string) (message string, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(log.Fields{
			"id": id,
			"err": fmt.Sprintf("%v", err),
			"took": time.Since(begin),
		}).Info("Disable job api")
	}(time.Now())

	message, err = mw.next.Disable(ctx, id)
	return
}

func (mw loggingMiddleware) Query(ctx context.Context, query job.Job) (results *[]job.Job, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(log.Fields{
			"query": fmt.Sprintf("%+v", query),
			"err": fmt.Sprintf("%v", err),
			"took": time.Since(begin),
		}).Info("Listing jobs api")
	}(time.Now())

	results, err = mw.next.List(ctx)
	return
}

func LoggingMiddleware(logger *log.Logger) Middleware {
	return func(next APISchedulerService) APISchedulerService {
		return loggingMiddleware{logger, next}
	}
}
