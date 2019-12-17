package engine

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
	"scheduleJob/common"
	"time"
)

type Executor struct {
	Engine Engine
}

var gExecutor *Executor

func NewExecutor(e Engine) *Executor {
	return &Executor{Engine: e}
}

func ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		err := gExecutor.Engine.Lock()
		if err != nil {
			log.Debugf("lock fail: %s", err.Error())
			return
		}
		defer gExecutor.Engine.UnLock()

		start := time.Now()

		log.Infof("job [%v] execute [%s]", info.Job.Name, info.Job.Command)

		// 执行任务
		cmd := exec.CommandContext(info.CancelCtx, "/bin/bash", "-c", info.Job.Command)

		// 执行输出
		output, err := cmd.CombinedOutput()

		end := time.Now()

		log.Infof("job [%v] execute [%s] output: %s, time: %v, err: %v",
			info.Job.Name, info.Job.Command, output, end.Sub(start), err)

		PushJobResult(&common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      output,
			Err:         err,
			StartTime:   start,
			EndTime:     end,
		})

	}()
}

func InitExecutor(e Engine) {
	gExecutor = NewExecutor(e)
}
