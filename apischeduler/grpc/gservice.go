package grpc

import (
	"strings"

	"google.golang.org/grpc/metadata"

	context "golang.org/x/net/context"

	"github.com/jinzhu/copier"
	"github.com/will7200/mjs/apischeduler/grpc/pb"
	"github.com/will7200/mjs/apischeduler/service"
	"github.com/will7200/mjs/job"
)

type wrapperService struct {
	service service.APISchedulerService
}

func NewGRPC(s service.APISchedulerService) (r pb.APISchedulerServer) {
	r = wrapperService{s}
	return r
}

func (this wrapperService) Add(ctx context.Context, arg *pb.AddRequest) (arp *pb.AddReply, e error) {
	reqjob := arg.Reqjob
	job := job.Job{}
	copier.Copy(&job, reqjob)
	md, _ := metadata.FromContext(ctx)
	for key, value := range md {
		ctx = context.WithValue(ctx, strings.ToUpper(key), value[0])
	}
	id, err := this.service.Add(ctx, job)
	if err != nil {
		return arp, err
	}
	arp = &pb.AddReply{Id: id}
	return arp, err
}
func (this wrapperService) Start(ctx context.Context, arg *pb.StartRequest) (*pb.StartReply, error) {
	return nil, nil
}
func (this wrapperService) Remove(ctx context.Context, arg *pb.RemoveRequest) (*pb.RemoveReply, error) {
	return nil, nil
}
func (this wrapperService) Change(ctx context.Context, arg *pb.ChangeRequest) (*pb.ChangeReply, error) {
	return nil, nil
}
func (this wrapperService) Get(ctx context.Context, arg *pb.GetRequest) (*pb.GetReply, error) {
	return nil, nil
}
func (this wrapperService) List(ctx context.Context, arg *pb.ListRequest) (*pb.ListReply, error) {
	return nil, nil
}
func (this wrapperService) Enable(ctx context.Context, arg *pb.EnableRequest) (*pb.EnableReply, error) {
	return nil, nil
}
func (this wrapperService) Disable(ctx context.Context, arg *pb.DisableRequest) (*pb.DisableReply, error) {
	return nil, nil
}

func (this wrapperService) Query(ctx context.Context, arg *pb.QueryRequest) (*pb.QueryReply, error) {
	jobquery := &job.Job{}
	copier.Copy(jobquery, arg.Query)
	list, err := this.service.Query(ctx, *jobquery)
	if err != nil {
		return &pb.QueryReply{}, err
	}
	qr := []*pb.Job{}
	for _, val := range *list {
		newjob := &pb.Job{}
		copier.Copy(newjob, val)
		qr = append(qr, newjob)
	}
	copier.Copy(qr, list)
	return &pb.QueryReply{Jobs: qr}, err
}
