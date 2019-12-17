package grpc_plugin

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"scheduleJob/common"
	"scheduleJob/config"
	"scheduleJob/model"
	proto "scheduleJob/plugin/grpc_plugin/proto"
	"scheduleJob/service"
)

type GrpcPlugin struct {
	Port         int
	JobService   service.JobService
	jobEventChan chan *common.JobEvent
	infos        map[string]*common.JobExecuteInfo
}

func (g *GrpcPlugin) Add(ctx context.Context, job *proto.JobReq) (*proto.Result, error) {
	log.Infof("Get Add Job Event. Name: %s, Crontab: %s, Command: %s", job.Name, job.Crontab, job.Command)
	err := g.JobService.Create(&model.Job{
		Name:    job.Name,
		Crontab: job.Crontab,
		Command: job.Command,
	})
	if err != nil {
		return &proto.Result{
			Code:    1,
			Message: err.Error(),
		}, err
	}
	g.PushJobEvent(common.BuildJobEvent(job.Name, job.Command, job.Crontab, common.SAVE_EVENT))
	return &proto.Result{Code: 0}, nil
}

func (g *GrpcPlugin) PushJobEvent(event *common.JobEvent) {
	g.jobEventChan <- event
}

func (g *GrpcPlugin) Delete(ctx context.Context, job *proto.JobReq) (*proto.Result, error) {
	log.Infof("Get Delete Job Event. Id: %d", job.Id)
	j, err := g.JobService.GetOne(job.Id)

	if err != nil {
		return &proto.Result{
			Code:    1,
			Message: err.Error(),
		}, err
	}

	log.Infof("Job Event. Name: %s, Crontab: %s, Command: %s", j.Name, j.Crontab, j.Command)

	err = g.JobService.Delete(job.Id)

	if err != nil {
		return &proto.Result{
			Code:    1,
			Message: err.Error(),
		}, err
	}

	g.PushJobEvent(common.BuildJobEvent(j.Name, j.Command, j.Crontab, common.DELETE_EVENT))

	return &proto.Result{Code: 0}, nil
}

func (g *GrpcPlugin) Update(ctx context.Context, job *proto.JobReq) (*proto.Result, error) {
	log.Infof("Get Update Job Event. Name: %s, Crontab: %s, Command: %s", job.Name, job.Crontab, job.Command)
	err := g.JobService.Update(&model.Job{
		Id:      job.Id,
		Name:    job.Name,
		Crontab: job.Crontab,
		Command: job.Command,
	})
	if err != nil {
		return &proto.Result{
			Code:    1,
			Message: err.Error(),
		}, err
	}
	g.PushJobEvent(common.BuildJobEvent(job.Name, job.Command, job.Crontab, common.SAVE_EVENT))

	return &proto.Result{Code: 0}, nil
}

func (g *GrpcPlugin) Start(context.Context, *proto.JobReq) (*proto.Result, error) {
	panic("implement me")
}

func (g *GrpcPlugin) Stop(context.Context, *proto.JobReq) (*proto.Result, error) {
	panic("implement me")
}

func (g *GrpcPlugin) GetRunning(ctx context.Context, job *proto.JobStatus) (*proto.JobList, error) {

	jobs := &proto.JobList{}

	for _, info := range g.infos {
		jobs.JobReq = append(jobs.JobReq, &proto.JobReq{
			Name:    info.Job.Name,
			Crontab: info.Job.CronExpr,
			Command: info.Job.Command,
		})
	}

	return jobs, nil
}

func NewGrpcPlugin(config config.Config, events chan *common.JobEvent, infos map[string]*common.JobExecuteInfo) *GrpcPlugin {
	return &GrpcPlugin{config.GrpcConfig.Port, service.NewJobService(), events, infos}
}

func (g *GrpcPlugin) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", g.Port))

	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	proto.RegisterJobMgrServer(grpcServer, g)

	log.Info("start grpc server ...")

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal(err)
	}
}

func (g *GrpcPlugin) Lock() error {
	return nil
}

func (g *GrpcPlugin) UnLock() error {
	return nil
}
