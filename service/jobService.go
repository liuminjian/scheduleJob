package service

import (
	"github.com/liuminjian/infra/base"
	"scheduleJob/dao"
	"scheduleJob/model"
)

type JobService interface {
	Create(job *model.Job) error
	GetAll() (*[]model.Job, error)
	Delete(id uint64) error
	GetOne(id uint64) (*model.Job, error)
	Update(job *model.Job) error
}

type jobService struct {
	jobDao *dao.JobDao
}

func NewJobService() JobService {
	return &jobService{}
}

func (j *jobService) initDao() {
	db := base.InitDBMaster()
	j.jobDao = dao.NewJobDao(db)
}

func (j *jobService) Create(job *model.Job) error {
	j.initDao()
	return j.jobDao.Create(job)
}

func (j *jobService) GetAll() (*[]model.Job, error) {
	j.initDao()
	return j.jobDao.GetAll()
}

func (j *jobService) Delete(id uint64) error {
	j.initDao()
	return j.jobDao.Delete(id)
}

func (j *jobService) GetOne(id uint64) (*model.Job, error) {
	j.initDao()
	return j.jobDao.GetOne(id)
}

func (j *jobService) Update(job *model.Job) error {
	j.initDao()
	return j.jobDao.Update(job)
}
