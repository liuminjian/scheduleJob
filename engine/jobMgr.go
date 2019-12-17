package engine

import (
	log "github.com/sirupsen/logrus"
	"scheduleJob/common"
	"scheduleJob/service"
)

type JobMgr struct {
	Engine Engine
}

func NewJobMgr(engine Engine) *JobMgr {
	return &JobMgr{Engine: engine}
}

func (j *JobMgr) Start() {
	j.InitJobs()
	go j.Engine.Run()
}

func (j *JobMgr) InitJobs() {
	jobService := service.NewJobService()
	jobs, err := jobService.GetAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, job := range *jobs {
		PushJobEvent(common.BuildJobEvent(job.Name, job.Command, job.Crontab, common.SAVE_EVENT))
	}
}

func InitJobMgr(engine Engine) {
	job := NewJobMgr(engine)
	job.Start()
}
