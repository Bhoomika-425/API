package service

import (
	"context"
	"errors"
	"project/internal/auth"
	"project/internal/cache"
	"project/internal/models"
	"project/internal/repository"
)

type Service struct {
	UserRepo repository.UserRepo
	auth     auth.UserAuth
	rdb      cache.Cache
}

//go:generate mockgen -source=ser.go -destination=mock-files/ser_mock.go -package=mock_files
type UserService interface {
	UserSignup(ctx context.Context, userData models.NewUser) (models.User, error)
	UserLogin(ctx context.Context, userData models.NewUser) (string, error)

	AddCompanyDetails(ctx context.Context, companyData models.Company) (models.Company, error)
	ViewAllCompanies(ctx context.Context) ([]models.Company, error)
	ViewCompanyDetails(ctx context.Context, cid uint64) (models.Company, error)

	ViewJob(ctx context.Context, cid uint64) ([]models.Jobs, error)
	AddJobDetails(ctx context.Context, jobRequest models.NewJobRequest, cid uint) (models.NewJobResponse, error)
	ViewAllJobs(ctx context.Context) ([]models.Jobs, error)
	ViewJobById(ctx context.Context, jid uint64) ([]models.Jobs, error)
	ApplyJobs(ctx context.Context, application []models.NewUserApplication) ([]models.NewUserApplication, error)
}

func NewService(userRepo repository.UserRepo, a auth.UserAuth, rdb cache.Cache) (UserService, error) {
	if userRepo == nil {
		return nil, errors.New("interface cannot be null")
	}
	return &Service{
		UserRepo: userRepo,
		auth:     a,
		rdb:      rdb,
	}, nil
}
