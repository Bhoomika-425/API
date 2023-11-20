package service

import (
	"context"
	"errors"
	"project/internal/auth"
	"project/internal/cache"
	"project/internal/models"
	"project/internal/repository"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestService_ViewJobById(t *testing.T) {
	type args struct {
		ctx context.Context
		jid uint64
	}
	tests := []struct {
		name string
		// s       *Service
		args             args
		want             []models.Jobs
		wantErr          bool
		mockRepoResponse func() ([]models.Jobs, error)
	}{
		{name: "error",
			args: args{
				ctx: context.Background(),
			},
			want:    []models.Jobs{},
			wantErr: false,
			mockRepoResponse: func() ([]models.Jobs, error) {
				return nil, errors.New("test error")
			},
		},
		{name: "success",
			args: args{
				ctx: context.Background(),
				jid: 15,
			},
			want: []models.Jobs{
				{Company: models.Company{
					Name:     "tcs",
					Location: "bang",
					Field:    "software",
				},
					Cid:          2,
					Name:         "developer",
					Salary:       "30000",
					NoticePeriod: "3 weeks",
				},
			},

			wantErr: false,
			mockRepoResponse: func() ([]models.Jobs, error) {
				return []models.Jobs{
					{
						Company: models.Company{
							Name:     "tcs",
							Location: "bang",
							Field:    "software",
						},
						Cid:          2,
						Name:         "developer",
						Salary:       "30000",
						NoticePeriod: "3 weeks",
					},
				}, nil
			},
		},
		{
			name: "invalid id",
			args: args{
				ctx: context.Background(),
				jid: 5,
			},
			want:    []models.Jobs{},
			wantErr: false,
			mockRepoResponse: func() ([]models.Jobs, error) {
				return []models.Jobs{}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().Jobbyjid(tt.args.ctx, tt.args.jid).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, &auth.Auth{}, cache.RDB{})
			got, err := s.ViewJobById(tt.args.ctx, tt.args.jid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ViewJobById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ViewJobById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ViewAllJobs(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		// s                *Service
		args             args
		want             []models.Jobs
		wantErr          bool
		mockRepoResponse func() ([]models.Jobs, error)
	}{
		{name: "error",
			want: nil,
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
			mockRepoResponse: func() ([]models.Jobs, error) {
				return nil, errors.New("test error")
			},
		},
		{
			name: "success",
			want: []models.Jobs{
				{
					Cid:          2,
					Name:         "tcs",
					Salary:       "30000",
					NoticePeriod: "3 weeks",
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
			mockRepoResponse: func() ([]models.Jobs, error) {
				return []models.Jobs{
					{
						Cid:          2,
						Name:         "tcs",
						Salary:       "30000",
						NoticePeriod: "3 weeks",
					},
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			mockCache := cache.NewMockCache(mc)
			mockRepo.EXPECT().FetchAllJobs(tt.args.ctx).Return(tt.mockRepoResponse()).AnyTimes()
			s, _ := NewService(mockRepo, &auth.Auth{}, mockCache)
			got, err := s.ViewAllJobs(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ViewAllJobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ViewAllJobs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ViewJob(t *testing.T) {
	type args struct {
		ctx context.Context
		cid uint64
	}
	tests := []struct {
		name string
		// s       *Service
		args             args
		want             []models.Jobs
		wantErr          bool
		mockRepoResponse func() ([]models.Jobs, error)
		mockRedis        func()
	}{
		{
			name: "success",
			want: []models.Jobs{
				{Cid: 2,
					Name:         "assosiate",
					Salary:       "50000",
					NoticePeriod: "3 days",
				},
			},
			args: args{
				ctx: context.Background(),
				cid: 4,
			},
			wantErr: false,
			mockRepoResponse: func() ([]models.Jobs, error) {
				return []models.Jobs{
					{
						Cid:          2,
						Name:         "assosiate",
						Salary:       "50000",
						NoticePeriod: "3 days",
					},
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
			mockRepoResponse: func() ([]models.Jobs, error) {
				return nil, errors.New("no jobs")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)

			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().Jobbycid(gomock.Any(), gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}

			s, _ := NewService(mockRepo, &auth.Auth{}, cache.RDB{})
			got, err := s.ViewJob(tt.args.ctx, tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ViewJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ViewJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_AddJobDetails(t *testing.T) {
	type args struct {
		ctx        context.Context
		jobRequest models.NewJobRequest
		cid        uint
	}
	tests := []struct {
		name string
		// s       *Service
		args             args
		want             models.NewJobResponse
		wantErr          bool
		mockRepoResponse func() (models.NewJobResponse, error)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				jobRequest: models.NewJobRequest{
					Company: models.Company{
						Name:     "ibm",
						Location: "bng",
						Field:    "sw",
					},
					Jid:          5,
					Name:         "assosiate",
					Salary:       "60000",
					NoticePeriod: "30 days",
					MinNp:        "10 days",
					MaxNP:        "26 days",
					Budget:       "1LPA",
					Description:  "this company is providing job opportunity",
					Minexp:       "2 years",
					MaxMax:       "5 years",
					Jobloc:       []uint{uint(1)},
					Skills:       []uint{uint(1)},
					Mode:         []uint{uint(1)},
					Degree:       []uint{uint(1)},
					Shift:        []uint{uint(1)},
					Type:         []uint{uint(1)},
				},
				cid: 5,
			},
			wantErr: false,
			want: models.NewJobResponse{
				ID: 2,
			},
			mockRepoResponse: func() (models.NewJobResponse, error) {
				return models.NewJobResponse{ID: 2}, nil
			},
		},
		{
			name: "failure",
			args: args{
				ctx: context.Background(),
				jobRequest: models.NewJobRequest{
					Company: models.Company{
						Name:     "ibm",
						Location: "bng",
						Field:    "sw",
					},
					Jid:          5,
					Name:         "assosiate",
					Salary:       "60000",
					NoticePeriod: "30 days",
					MinNp:        "10 days",
					MaxNP:        "26 days",
					Budget:       "1LPA",
					Description:  "this company is providing job opportunity",
					Minexp:       "2 years",
					MaxMax:       "5 years",
					Jobloc:       []uint{uint(1)},
					Skills:       []uint{uint(1)},
					Mode:         []uint{uint(1)},
					Degree:       []uint{uint(1)},
					Shift:        []uint{uint(1)},
					Type:         []uint{uint(1)},
				},
				cid: 5,
			},
			wantErr: true,
			want:    models.NewJobResponse{},
			mockRepoResponse: func() (models.NewJobResponse, error) {
				return models.NewJobResponse{}, errors.New("error in creating jobs")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			mockCache := cache.NewMockCache(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().CreateUserJob(gomock.Any(), gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, &auth.Auth{}, mockCache)
			got, err := s.AddJobDetails(tt.args.ctx, tt.args.jobRequest, tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AddJobDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.AddJobDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestService_ApplyJobs(t *testing.T) {
// 	type args struct {
// 		ctx          context.Context
// 		applications []models.NewUserApplication
// 	}
// 	tests := []struct {
// 		name string
// 		// s       *Service
// 		args             args
// 		want             []models.NewUserApplication
// 		wantErr          bool
// 		mockRepoResponse func() (models.Jobs, error)
// 	}{
// 		{
// 			name: "error in mock fun",
// 			args: args{ctx: context.Background(),
// 				applications: []models.NewUserApplication{
// 					{Name: "Jhon",
// 						Age: "25",
// 						ID:  1,
// 						Jobs: models.Requestfield{
// 							Name:         "assosiate",
// 							NoticePeriod: 3,
// 							Experience:   2,
// 							Locations: []uint{
// 								uint(1), uint(2),
// 							},
// 							TechnologyStacks: []uint{
// 								uint(1), uint(2),
// 							},
// 							Degree: []uint{
// 								uint(1), uint(2),
// 							},
// 							Shifts: []uint{
// 								uint(1), uint(2),
// 							},
// 						},
// 					},
// 				},
// 			},
// 			want:    nil,
// 			wantErr: false,
// 			mockRepoResponse: func() (models.Jobs, error) {
// 				return models.Jobs{}, errors.New("error from mock function")
// 			},
// 		},
// 		{
// 			name: "experience does not match",
// 			args: args{
// 				ctx: context.Background(),
// 				applications: []models.NewUserApplication{
// 					{
// 						Name: "Jhon",
// 						Age:  "25",
// 						ID:   1,
// 						Jobs: models.Requestfield{
// 							Name:         "assosiate",
// 							NoticePeriod: 3,
// 							Experience:   5,
// 							Locations: []uint{
// 								uint(1), uint(2),
// 							},
// 							// 	TechnologyStacks: []uint{
// 							// 		uint(1), uint(2),
// 							// 	},
// 							// 	Degree: []uint{
// 							// 		uint(1), uint(2),
// 							// 	},
// 							// 	Shifts: []uint{
// 							// 		uint(1), uint(2),
// 							// 	},
// 						},
// 					},
// 				},
// 			},
// 			want:    nil,
// 			wantErr: false,
// 			mockRepoResponse: func() (models.Jobs, error) {
// 				return models.Jobs{
// 					Model: gorm.Model{
// 						ID: 1,
// 					},
// 					Company: models.Company{
// 						Model: gorm.Model{
// 							ID: 1,
// 						},
// 					},
// 					Cid:          1,
// 					Name:         "assosiate",
// 					Salary:       "80000",
// 					NoticePeriod: "5",
// 					MinNp:        "3",
// 					MaxNP:        "4",
// 					Budget:       "15lpa",
// 					Description:  "this is the software job",
// 					Minexp:       "0",
// 					MaxMax:       "10",
// 					// Locations: []models.Location{
// 					// 	{
// 					// 		Model: gorm.Model{
// 					// 			ID: 1,
// 					// 		},
// 					// 	},
// 					// 	{
// 					// 		Model: gorm.Model{
// 					// 			ID: 2,
// 					// 		},
// 					// 	},
// 					// },
// 					// TechnologyStacks: []models.TechnologyStack{
// 					// 	{
// 					// 		Model: gorm.Model{
// 					// 			ID: 1,
// 					// 		},
// 					// 	},
// 					// 	{
// 					// 		Model: gorm.Model{
// 					// 			ID: 2,
// 					// 		},
// 					// 	},
// 					// },
// 					//  WorkModes: []models.WorkMode{
// 					//  	{
// 					//  		Model: gorm.Model{
// 					// 			ID: 1,
// 					// 		},
// 					// 	},
// 					// 	{
// 					// 		Model: gorm.Model{
// 					// 			ID: 2,
// 					// 		},
// 					// 	},
// 					// },
// 					// Qualifications: []models.Qualification{
// 					// 	{
// 					// 		Model: gorm.Model{
// 					// 			ID: 1,
// 					// 		},
// 					// 	},
// 					// 	{
// 					// 		Model: gorm.Model{
// 					// 			ID: 2,
// 					// 		},
// 					// 	},
// 					// },
// 				}, nil
// 			},
// 		},
// 		{
// 			name: "error in fetching location",
// 			args: args{
// 				ctx: context.Background(),
// 				applications: []models.NewUserApplication{
// 					{
// 						Name: "Jhon",
// 						Age:  "25",
// 						ID:   1,
// 						Jobs: models.Requestfield{
// 							Name:         "assosiate",
// 							NoticePeriod: 3,
// 							Experience:   2,
// 							Locations: []uint{
// 								uint(1), uint(2),
// 							},
// 							// TechnologyStacks: []uint{
// 							// 	uint(1), uint(2),
// 							// },
// 							// Degree: []uint{
// 							// 	uint(1), uint(2),
// 							// },
// 							// Shifts: []uint{
// 							// 	uint(1), uint(2),
// 							// },
// 						},
// 					},
// 				},
// 			},
// 			want:    nil,
// 			wantErr: false,
// 			mockRepoResponse: func() (models.Jobs, error) {
// 				return models.Jobs{
// 					Model: gorm.Model{
// 						ID: 1,
// 					},
// 					Company: models.Company{
// 						Model: gorm.Model{
// 							ID: 1,
// 						},
// 					},
// 					Cid:          1,
// 					Name:         "assosiate",
// 					Salary:       "80000",
// 					NoticePeriod: "5",
// 					MinNp:        "3",
// 					MaxNP:        "4",
// 					Budget:       "15lpa",
// 					Description:  "this is the software job",
// 					Minexp:       "0",
// 					MaxMax:       "10",
// 					Locations: []models.Location{
// 						{
// 							Model: gorm.Model{
// 								ID: 1,
// 							},
// 						},
// 						{
// 							Model: gorm.Model{
// 								ID: 3,
// 							},
// 						},
// 					},
// 					TechnologyStacks: []models.TechnologyStack{
// 						{
// 							Model: gorm.Model{
// 								ID: 1,
// 							},
// 						},
// 						{
// 							Model: gorm.Model{
// 								ID: 2,
// 							},
// 						},
// 					},
// 				}, nil
// 			},
// 		},
// 		{
// 			name: "error in checking technology stack",
// 			args: args{
// 				ctx: context.Background(),
// 				applications: []models.NewUserApplication{
// 					{
// 						Name: "Jhon",
// 						Age:  "25",
// 						ID:   1,
// 						Jobs: models.Requestfield{
// 							Name:         "assosiate",
// 							NoticePeriod: 3,
// 							Experience:   2,
// 							Locations: []uint{
// 								uint(1), uint(2),
// 							},
// 							TechnologyStacks: []uint{
// 								uint(1), uint(2),
// 							},
// 							Degree: []uint{
// 								uint(1), uint(2),
// 							},
// 						},
// 					},
// 				},
// 			},
// 			want:    nil,
// 			wantErr: false,
// 			mockRepoResponse: func() (models.Jobs, error) {
// 				return models.Jobs{
// 					Model: gorm.Model{
// 						ID: 1,
// 					},
// 					Company: models.Company{
// 						Model: gorm.Model{
// 							ID: 1,
// 						},
// 					},
// 					Cid:          1,
// 					Name:         "assosiate",
// 					Salary:       "80000",
// 					NoticePeriod: "5",
// 					MinNp:        "3",
// 					MaxNP:        "4",
// 					Budget:       "15lpa",
// 					Description:  "this is the software job",
// 					Minexp:       "0",
// 					MaxMax:       "10",
// 					TechnologyStacks: []models.TechnologyStack{
// 						{
// 							Model: gorm.Model{
// 								ID: 0,
// 							},
// 						},
// 						{
// 							Model: gorm.Model{
// 								ID: 5,
// 							},
// 						},
// 					},
// 				}, nil
// 			},
// 		},
// 		// {
// 		// 	name: "error in checking qualifications",
// 		// 	args: args{
// 		// 		ctx: context.Background(),
// 		// 		applications: []models.NewUserApplication{
// 		// 			{
// 		// 				Name: "Jhon",
// 		// 				Age:  "25",
// 		// 				ID:   1,
// 		// 				Jobs: models.Requestfield{
// 		// 					Name:         "assosiate",
// 		// 					NoticePeriod: 3,
// 		// 					Experience:   2,
// 		// 					Locations: []uint{
// 		// 						uint(1), uint(2),
// 		// 					},
// 		// 					TechnologyStacks: []uint{
// 		// 						uint(1), uint(2),
// 		// 					},
// 		// 					Degree: []uint{
// 		// 						uint(1), uint(2),
// 		// 					},
// 		// 				},
// 		// 			},
// 		// 		},
// 		// 	},
// 		// 	want:    nil,
// 		// 	wantErr: false,
// 		// 	mockRepoResponse: func() (models.Jobs, error) {
// 		// 		return models.Jobs{
// 		// 			Model: gorm.Model{
// 		// 				ID: 1,
// 		// 			},
// 		// 			Company: models.Company{
// 		// 				Model: gorm.Model{
// 		// 					ID: 1,
// 		// 				},
// 		// 			},
// 		// 			Cid:          1,
// 		// 			Name:         "assosiate",
// 		// 			Salary:       "80000",
// 		// 			NoticePeriod: "5 weeks",
// 		// 			MinNp:        "5days",
// 		// 			MaxNP:        "4weeks",
// 		// 			Budget:       "15lpa",
// 		// 			Description:  "this is the software job",
// 		// 			Minexp:       "1year",
// 		// 			MaxMax:       "2years",
// 		// 			WorkModes: []models.WorkMode{
// 		// 				{
// 		// 					Model: gorm.Model{
// 		// 						ID: 1,
// 		// 					},
// 		// 				},
// 		// 				{
// 		// 					Model: gorm.Model{
// 		// 						ID: 2,
// 		// 					},
// 		// 				},
// 		// 			},
// 		// 			Qualifications: []models.Qualification{
// 		// 				{
// 		// 					Model: gorm.Model{
// 		// 						ID: 0,
// 		// 					},
// 		// 				},
// 		// 				{
// 		// 					Model: gorm.Model{
// 		// 						ID: 2,
// 		// 					},
// 		// 				},
// 		// 			},
// 		// 		}, nil
// 		// 	},
// 		// },

// 		{
// 			name: "error in getting shifts",
// 			args: args{
// 				ctx: context.Background(),
// 				applications: []models.NewUserApplication{
// 					{
// 						Name: "Jhon",
// 						Age:  "25",
// 						ID:   1,
// 						Jobs: models.Requestfield{
// 							Name:         "assosiate",
// 							NoticePeriod: 3,
// 							Experience:   2,
// 							// Locations: []uint{
// 							// 	uint(1), uint(2),
// 							// },
// 							// TechnologyStacks: []uint{
// 							// 	uint(1), uint(2),
// 							// },
// 							Shifts: []uint{
// 								uint(1), uint(2),
// 							},
// 						},
// 					},
// 				},
// 			},
// 			want:    nil,
// 			wantErr: false,
// 			mockRepoResponse: func() (models.Jobs, error) {
// 				return models.Jobs{
// 					Model: gorm.Model{
// 						ID: 1,
// 					},
// 					Company: models.Company{
// 						Model: gorm.Model{
// 							ID: 1,
// 						},
// 					},
// 					Cid:          1,
// 					Name:         "assosiate",
// 					Salary:       "80000",
// 					NoticePeriod: "5 weeks",
// 					MinNp:        "5days",
// 					MaxNP:        "4weeks",
// 					Budget:       "15lpa",
// 					Description:  "this is the software job",
// 					Minexp:       "1year",
// 					MaxMax:       "2years",
// 					Locations: []models.Location{
// 						{
// 							Model: gorm.Model{
// 								ID: 1,
// 							},
// 						},
// 						{
// 							Model: gorm.Model{
// 								ID: 2,
// 							},
// 						},
// 					},
// 					// TechnologyStacks: []models.TechnologyStack{
// 					// 	{
// 					// 		Model: gorm.Model{
// 					// 			ID: 1,
// 					// 		},
// 					// 	},
// 					// 	{
// 					// 		Model: gorm.Model{
// 					// 			ID: 2,
// 					// 		},
// 					// 	},
// 					// },
// 					Shifts: []models.Shift{
// 						{
// 							Model: gorm.Model{
// 								ID: 0,
// 							},
// 						},
// 						{
// 							Model: gorm.Model{
// 								ID: 1,
// 							},
// 						},
// 					},
// 				}, nil
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mc := gomock.NewController(t)
// 			mockRepo := repository.NewMockUserRepo(mc)
// 			if tt.mockRepoResponse != nil {
// 				mockRepo.EXPECT().CreateApplication(gomock.Any(), gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
// 			}
// 			for _, v := range tt.args.applications {
// 				if v.ID == uint(1) {
// 					mockRepo.EXPECT().CreateApplication(gomock.Any(), v.ID).Return(models.Jobs{
// 						Model: gorm.Model{
// 							ID: 1,
// 						},
// 						Company: models.Company{
// 							Model: gorm.Model{
// 								ID: 1,
// 							},
// 						},
// 						Cid:          1,
// 						Name:         "assosiate",
// 						Salary:       "80000",
// 						NoticePeriod: "5 weeks",
// 						MinNp:        "5days",
// 						MaxNP:        "4weeks",
// 						Budget:       "15lpa",
// 						Description:  "this is the software job",
// 						Minexp:       "1year",
// 						MaxMax:       "2years",
// 						Locations: []models.Location{
// 							{
// 								Model: gorm.Model{
// 									ID: 1,
// 								},
// 							},
// 							{
// 								Model: gorm.Model{
// 									ID: 2,
// 								},
// 							},
// 						},
// 						TechnologyStacks: []models.TechnologyStack{
// 							{
// 								Model: gorm.Model{
// 									ID: 1,
// 								},
// 							},
// 							{
// 								Model: gorm.Model{
// 									ID: 2,
// 								},
// 							},
// 						},
// 						WorkModes: []models.WorkMode{
// 							{
// 								Model: gorm.Model{
// 									ID: 1,
// 								},
// 							},
// 							{
// 								Model: gorm.Model{
// 									ID: 2,
// 								},
// 							},
// 						},
// 						Qualifications: []models.Qualification{
// 							{
// 								Model: gorm.Model{
// 									ID: 1,
// 								},
// 							},
// 							{
// 								Model: gorm.Model{
// 									ID: 2,
// 								},
// 							},
// 						},
// 						Shifts: []models.Shift{
// 							{
// 								Model: gorm.Model{
// 									ID: 1,
// 								},
// 							},
// 							{
// 								Model: gorm.Model{
// 									ID: 2,
// 								},
// 							},
// 						},
// 						Jobtypes: []models.Jobtype{
// 							{
// 								Model: gorm.Model{
// 									ID: 1,
// 								},
// 							},
// 							{
// 								Model: gorm.Model{
// 									ID: 2,
// 								},
// 							},
// 						},
// 					}, nil).AnyTimes()
// 				}
// 			}
// 			s, err := NewService(mockRepo, &auth.Auth{}, cache.RDB{})
// 			if err != nil {
// 				t.Errorf("error is initializing the repo layer")
// 				return
// 			}
// 			got, err := s.ApplyJobs(tt.args.ctx, tt.args.applications)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Service.ApplyJobs() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Service.ApplyJobs() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
