package common

import (
	"context"
	"github.com/gorhill/cronexpr"
	"time"
)

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

type EventType int

const (
	SAVE_EVENT   EventType = 1
	DELETE_EVENT EventType = 2
	START_EVENT  EventType = 3
	STOP_EVENT   EventType = 4
)

type JobEvent struct {
	Job *Job
	EventType
}

type SchedulePlan struct {
	Job      *Job
	Expr     *cronexpr.Expression
	NextTime time.Time
}

type JobExecuteInfo struct {
	Job        *Job
	PlainTime  time.Time
	RealTime   time.Time
	CancelCtx  context.Context
	CancelFunc context.CancelFunc
}

type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo
	Output      []byte
	Err         error
	StartTime   time.Time
	EndTime     time.Time
}

func BuildSchedulePlan(event *JobEvent) (*SchedulePlan, error) {

	expr, err := cronexpr.Parse(event.Job.CronExpr)

	if err != nil {
		return nil, err
	}

	return &SchedulePlan{
		Job:      event.Job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}, nil
}

func BuildExecuteInfo(plain *SchedulePlan) *JobExecuteInfo {
	ctx, cancel := context.WithCancel(context.Background())
	return &JobExecuteInfo{
		Job:        plain.Job,
		PlainTime:  plain.NextTime,
		RealTime:   time.Now(),
		CancelCtx:  ctx,
		CancelFunc: cancel,
	}
}

func BuildJobEvent(name string, command string, crontab string, event EventType) *JobEvent {
	return &JobEvent{
		Job: &Job{
			Name:     name,
			Command:  command,
			CronExpr: crontab,
		},
		EventType: event,
	}
}
