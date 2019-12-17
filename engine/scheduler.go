package engine

import (
	log "github.com/sirupsen/logrus"
	"scheduleJob/common"
	"time"
)

var gScheduler *scheduler

type scheduler struct {
	jobEventChan    chan *common.JobEvent
	jobPlainTable   map[string]*common.SchedulePlan
	jobExecuteTable map[string]*common.JobExecuteInfo
	jobResultChan   chan *common.JobExecuteResult
}

func NewScheduler() *scheduler {
	return &scheduler{
		jobEventChan:    make(chan *common.JobEvent),
		jobPlainTable:   make(map[string]*common.SchedulePlan),
		jobExecuteTable: make(map[string]*common.JobExecuteInfo),
		jobResultChan:   make(chan *common.JobExecuteResult),
	}
}

func (s *scheduler) scheduleLoop() {

	schedulerAfter := s.trySchedule()

	// 等到最近的任务
	scheduleTimer := time.NewTimer(schedulerAfter)

	// 定时job

	for {
		select {
		case jobEvent := <-s.jobEventChan:
			s.handleJobEvent(jobEvent)
		case <-scheduleTimer.C:
		case jobResult := <-s.jobResultChan:
			s.handleJobResult(jobResult)
		}
		schedulerAfter = s.trySchedule()
		scheduleTimer.Reset(schedulerAfter)
	}

}

func (s *scheduler) handleJobEvent(event *common.JobEvent) {
	switch event.EventType {
	case common.SAVE_EVENT:
		plain, err := common.BuildSchedulePlan(event)
		if err != nil {
			log.Error(err)
			return
		}
		s.jobPlainTable[event.Job.Name] = plain
	case common.DELETE_EVENT:
		if _, exist := s.jobPlainTable[event.Job.Name]; exist {
			delete(s.jobPlainTable, event.Job.Name)
		}
	case common.START_EVENT:
	case common.STOP_EVENT:
	}
}

func (s *scheduler) handleJobResult(result *common.JobExecuteResult) {
	delete(s.jobExecuteTable, result.ExecuteInfo.Job.Name)
}

func (s *scheduler) tryStartJob(plan *common.SchedulePlan) {
	if _, exist := s.jobExecuteTable[plan.Job.Name]; exist {
		return
	}

	executeInfo := common.BuildExecuteInfo(plan)

	s.jobExecuteTable[plan.Job.Name] = executeInfo

	// 执行任务
	ExecuteJob(executeInfo)
}

func (s *scheduler) trySchedule() (schedulerAfter time.Duration) {
	if len(s.jobPlainTable) == 0 {
		schedulerAfter = 1 * time.Second
		return
	}

	now := time.Now()

	var nearTime *time.Time

	for _, jobPlan := range s.jobPlainTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			s.tryStartJob(jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(now)
		}

		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}

	if nearTime == nil {
		schedulerAfter = 1 * time.Second
		return
	}

	schedulerAfter = (*nearTime).Sub(now)
	return
}

func PushJobEvent(event *common.JobEvent) {
	gScheduler.jobEventChan <- event
}

func PushJobResult(result *common.JobExecuteResult) {
	gScheduler.jobResultChan <- result
}

func (s *scheduler) Start() {
	go s.scheduleLoop()
}

func GetEventChan() chan *common.JobEvent {
	return gScheduler.jobEventChan
}

func GetJobExecuteTable() map[string]*common.JobExecuteInfo {
	return gScheduler.jobExecuteTable
}

func InitScheduler() {
	gScheduler = NewScheduler()
	gScheduler.Start()
}
