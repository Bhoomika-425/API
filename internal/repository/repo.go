package repository

import (
	"context"
	"errors"
	"project/internal/models"

	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

//go:generate mockgen -source=repo.go -destination=repo_mock.go -package=repository
type UserRepo interface {
	CreateUser(ctx context.Context, userData models.User) (models.User, error)
	Userbyemail(ctx context.Context, email string) (models.User, error)

	CreateUserCompany(ctx context.Context, companyData models.Company) (models.Company, error)
	Companies(ctx context.Context) ([]models.Company, error)
	CompanyById(ctx context.Context, cid uint64) (models.Company, error)

	CreateUserJob(ctx context.Context, jobData models.Jobs) (models.NewJobResponse, error)
	Jobbycid(ctx context.Context, cid uint64) ([]models.Jobs, error)
	FetchAllJobs(ctx context.Context) ([]models.Jobs, error)
	Jobbyjid(ctx context.Context, jid uint64) (models.Jobs, error)
	CreateApplication(ctx context.Context, jid uint) (models.Jobs, error)
}

func NewRepository(db *gorm.DB) (UserRepo, error) {
	if db == nil {
		return nil, errors.New("db cannot be null")
	}
	return &Repo{
		DB: db,
	}, nil
}
