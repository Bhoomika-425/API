package repository

import (
	"context"
	"errors"
	"project/internal/models"

	"github.com/rs/zerolog/log"
)

func (r *Repo) CreateUserCompany(ctx context.Context, companyData models.Company) (models.Company, error) {
	result := r.DB.Create(&companyData)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.Company{}, errors.New("could not create the company")
	}
	return companyData, nil
}

func (r *Repo) Companies(ctx context.Context) ([]models.Company, error) {
	var userDetails []models.Company
	result := r.DB.Find(&userDetails)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return nil, errors.New("could not find the companies")
	}
	return userDetails, nil
}

func (r *Repo) CompanyById(ctx context.Context, cid uint64) (models.Company, error) {
	var companyData models.Company
	result := r.DB.Where("id = ?", cid).First(&companyData)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.Company{}, errors.New("could not find the company")
	}
	return companyData, nil
}
