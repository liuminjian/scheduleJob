package engine

import (
	"fmt"
	"scheduleJob/common"
	"scheduleJob/config"
	"scheduleJob/plugin/grpc_plugin"
)

type Engine interface {
	Run()
	Lock() error
	UnLock() error
}

func NewEngine(config config.Config, events chan *common.JobEvent, table map[string]*common.JobExecuteInfo) Engine {
	if config.Engine == "grpc" {
		return grpc_plugin.NewGrpcPlugin(config, events, table)
	}

	panic(fmt.Sprintf("no suitable engine for %s", config.Engine))
}
