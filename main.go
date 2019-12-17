package main

import (
	"flag"
	"github.com/liuminjian/infra/base"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"scheduleJob/config"
	"scheduleJob/engine"
	"scheduleJob/store"
	"time"
)

var (
	configFile string
)

func iniArgs() {
	flag.StringVar(&configFile, "config", "./config/conf.yaml", "yaml config file")
	flag.Parse()
}

func GetConfig() config.Config {
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("read yaml file fail:%s", err.Error())
	}
	var cfg config.Config
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Unmarshal yaml file fail:%s", err.Error())
	}

	if err := base.ValidateStruct(cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}

func SetLog(logPath string) {
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Fatalf("Failed to open log file: %s", err.Error())
	}
}

func main() {

	// 初始化参数
	iniArgs()

	// 初始化校验
	base.InitValidator()

	// 读配置文件
	cfg := GetConfig()

	// 初始化日志
	level, _ := log.ParseLevel(cfg.LogConfig.Level)
	hour := time.Hour.Nanoseconds()
	base.InitLog(cfg.LogConfig.Path, level, time.Duration(cfg.LogConfig.MaxAge*hour),
		time.Duration(cfg.LogConfig.RotationTime*hour))

	// 初始化数据库
	store.InitStore(cfg)

	// 初始化调度
	engine.InitScheduler()

	// 初始化插件
	e := engine.NewEngine(cfg, engine.GetEventChan(), engine.GetJobExecuteTable())

	// 初始化任务管理
	engine.InitJobMgr(e)

	// 初始化执行器
	engine.InitExecutor(e)

	// 等待
	for {
		time.Sleep(1 * time.Second)
	}
}
