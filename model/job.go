package model

import "time"

type Job struct {
	Id         uint64    `gorm:"primary_key" json:"id"`
	Name       string    `gorm:"not null;type varchar(255);unique" sql:"comment:'任务名称'" json:"name"`
	Crontab    string    `gorm:"not null;type varchar(255)" sql:"comment:'定时策略'" json:"crontab"`
	ExpireTime time.Time `gorm:"" sql:"comment:'有效期'" json:"expireTime"`
	Command    string    `gorm:"not null;type varchar(255)" sql:"comment:'任务命令'" json:"command"`
}
