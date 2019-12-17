package store

import (
	"github.com/liuminjian/infra/base"
	"scheduleJob/config"
)

func InitStore(cfg config.Config) {

	if cfg.StoreType == "mysql" {
		mysqlConfig := &base.MysqlConfig{
			Host:     cfg.MysqlConfig.Ip,
			Port:     cfg.MysqlConfig.Port,
			User:     cfg.MysqlConfig.User,
			Password: cfg.MysqlConfig.Password,
			Database: cfg.MysqlConfig.Database,
			Debug:    cfg.MysqlConfig.Debug,
		}
		base.InitDBMaster(mysqlConfig)
	}

}
