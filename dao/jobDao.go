package dao

import (
	"github.com/jinzhu/gorm"
	"scheduleJob/model"
)

type JobDao struct {
	db *gorm.DB
}

func NewJobDao(db *gorm.DB) *JobDao {
	return &JobDao{db: db}
}

func (j *JobDao) Create(job *model.Job) error {
	j.db.Create(job)
	return j.db.Error
}

func (j *JobDao) Delete(id uint64) error {
	j.db.Delete(&model.Job{
		Id: id,
	})
	return j.db.Error
}

func (j *JobDao) GetAll() (*[]model.Job, error) {
	jobs := make([]model.Job, 0)
	j.db.Find(&jobs)
	return &jobs, j.db.Error
}

func (j *JobDao) GetOne(id uint64) (*model.Job, error) {
	job := &model.Job{Id: id}
	j.db.Where(job).First(job)
	return job, j.db.Error
}

func (j *JobDao) Update(job *model.Job) error {
	j.db.Model(job).Updates(*job)
	return j.db.Error
}
