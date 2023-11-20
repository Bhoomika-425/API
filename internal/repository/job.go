package repository

import (
	"context"
	"errors"
	"project/internal/models"

	"github.com/rs/zerolog/log"
)

func (r *Repo) Jobbyjid(ctx context.Context, jid uint64) ([]models.Jobs, error) {
	var jobData []models.Jobs
	result := r.DB.Where("id = ?", jid).Find(&jobData)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return nil, errors.New("could not create the jobs")
	}
	return jobData, nil
}

func (r *Repo) CreateUserJob(ctx context.Context, jobData models.Jobs) (models.NewJobResponse, error) {

	result := r.DB.Create(&jobData)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.NewJobResponse{}, errors.New("could not create the jobs")
	}
	return models.NewJobResponse{ID: jobData.ID}, nil
}

func (r *Repo) FetchAllJobs(ctx context.Context) ([]models.Jobs, error) {
	var jobDatas []models.Jobs
	result := r.DB.Find(&jobDatas)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return nil, errors.New("could not find the jobs")
	}
	return jobDatas, nil
}

func (r *Repo) Jobbycid(ctx context.Context, cid uint64) ([]models.Jobs, error) {
	var jobData []models.Jobs
	result := r.DB.Where("cid = ?", cid).Find(&jobData)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return []models.Jobs{}, errors.New("could not find the company")
	}
	return jobData, nil
}
func (r *Repo) CreateApplication(ctx context.Context, jid uint) (models.Jobs, error) {
	var jobData models.Jobs

	// Preload related data using GORM's Preload method
	result := r.DB.Preload("Company").
		Preload("Locations").
		Preload("TechnologyStacks").
		Preload("Qualifications").
		Preload("Shifts").
		Where("id = ?", jid).
		Find(&jobData)

	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.Jobs{}, result.Error
	}

	return jobData, nil
}
