package service

import (
	"context"
	"errors"
	"project/internal/auth"
	"project/internal/models"
	"project/internal/repository"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestService_AddCompanyDetails(t *testing.T) {
	type args struct {
		ctx         context.Context
		companyData models.Company
	}
	tests := []struct {
		name string
		// s       *Service
		args             args
		want             models.Company
		wantErr          bool
		mockRepoResponse func() (models.Company, error)
	}{
		{
			name: "success",
			args: args{
				ctx:         context.Background(),
				companyData: models.Company{},
			},
			want: models.Company{
				Name:     "ibm",
				Location: "hyderabad",
				Field:    "software",
			},
			wantErr: false,
			mockRepoResponse: func() (models.Company, error) {
				return models.Company{
					Name:     "ibm",
					Location: "hyderabad",
					Field:    "software",
				}, nil
			},
		},
		{name: "failure",
			args: args{
				ctx:         context.Background(),
				companyData: models.Company{},
			},
			want:    models.Company{},
			wantErr: true,
			mockRepoResponse: func() (models.Company, error) {
				return models.Company{}, errors.New("fields are required")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().CreateUserCompany(gomock.Any(), gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, &auth.Auth{})
			got, err := s.AddCompanyDetails(tt.args.ctx, tt.args.companyData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AddCompanyDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.AddCompanyDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ViewAllCompanies(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		// s       *Service
		args             args
		want             []models.Company
		wantErr          bool
		mockRepoResponse func() ([]models.Company, error)
	}{
		{name: "success",
			args: args{
				ctx: context.Background(),
			},
			want: []models.Company{{
				Name:     "infosys",
				Location: "bng",
				Field:    "hardware",
			},
			},
			wantErr: false,
			mockRepoResponse: func() ([]models.Company, error) {
				return []models.Company{{
					Name:     "infosys",
					Location: "bng",
					Field:    "hardware"},
				}, nil
			},
		},
		{
			name: "failure",
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
			mockRepoResponse: func() ([]models.Company, error) {
				return nil, errors.New("no company")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().Companies(tt.args.ctx).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, &auth.Auth{})
			got, err := s.ViewAllCompanies(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ViewAllCompanies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ViewAllCompanies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ViewCompanyDetails(t *testing.T) {
	type args struct {
		ctx context.Context
		cid uint64
	}
	tests := []struct {
		name string
		// s       *Service
		args             args
		want             models.Company
		wantErr          bool
		mockRepoResponse func() (models.Company, error)
	}{
		{
			name: "failure",
			args: args{
				ctx: context.Background(),
				cid: 0,
			},
			want:    models.Company{},
			wantErr: true,
			mockRepoResponse: func() (models.Company, error) {
				return models.Company{}, errors.New("id should not be zero")
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				cid: 9,
			},
			want: models.Company{
				Name:     "microsoft",
				Location: "bng",
				Field:    "sofware",
			},
			wantErr: false,
			mockRepoResponse: func() (models.Company, error) {
				return models.Company{
					Name:     "microsoft",
					Location: "bng",
					Field:    "sofware",
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().CompanyById(tt.args.ctx, tt.args.cid).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, &auth.Auth{})
			got, err := s.ViewCompanyDetails(tt.args.ctx, tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ViewCompanyDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ViewCompanyDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}
