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

func TestService_UserLogin(t *testing.T) {
	type args struct {
		ctx      context.Context
		userData models.NewUser
	}
	tests := []struct {
		name string
		// s       Service
		args             args
		want             string
		wantErr          bool
		mockRepoResponse func() (models.User, error)
	}{
		{
			name: "error in mail",
			args: args{
				ctx: context.Background(),
				userData: models.NewUser{
					Email:    "abc@gmail.com",
					Password: "bhoomi25",
				},
			},
			want:    "",
			wantErr: true,
			mockRepoResponse: func() (models.User, error) {
				return models.User{}, errors.New("error")
			},
		},
		{
			name: "error in hashing",
			args: args{
				ctx: context.Background(),
				userData: models.NewUser{
					Email:    "abc@gmail.com",
					Password: "bhoom#$@@#",
				},
			},
			want:    "",
			wantErr: true,
			mockRepoResponse: func() (models.User, error) {
				return models.User{}, errors.New("error")
			},
		},
		// 	{name:"success",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		userData: models.NewUser{
		// 			Email:    "abc@gmail.com",
		// 			Password: "bhoom#$@@#",
		// 		},
		// 	},
		// 	want:    "bhoom#$@@#",
		// 	wantErr: false,
		// 	mockRepoResponse: func() (models.User, error) {
		// 		return "$2y$10$zrRsu29vEQ2MqxKY2BEHC.D0irqpcYoNh/RusMRR16TiJQ5IaG0b6",
		// 	},nil
		// },

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().Userbyemail(tt.args.ctx, tt.args.userData.Email).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, &auth.Auth{})
			got, err := s.UserLogin(tt.args.ctx, tt.args.userData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UserLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.UserLogin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_UserSignup(t *testing.T) {
	type args struct {
		ctx      context.Context
		userData models.NewUser
	}
	tests := []struct {
		name string
		// s       Service
		args             args
		want             models.User
		wantErr          bool
		mockRepoResponse func() (models.User, error)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				userData: models.NewUser{
					Username: "abc",
					Email:    "abc@gmail.com",
					Password: "990",
				},
			},
			want: models.User{
				Username:     "abc",
				Email:        "abc@gmail.com",
				PasswordHash: "9980",
			},
			wantErr: false,
			mockRepoResponse: func() (models.User, error) {
				return models.User{
					Username:     "abc",
					Email:        "abc@gmail.com",
					PasswordHash: "9980",
				}, nil
			},
		},

		{
			name: "email notcorrect",
			args: args{
				ctx: context.Background(),
				userData: models.NewUser{
					Username: "abc",
					Email:    "abc@mail.com",
					Password: "990",
				},
			},
			want:    models.User{},
			wantErr: true,
			mockRepoResponse: func() (models.User, error) {
				return models.User{}, errors.New("email is wrong")
			},
		},
		{
			name: "invalid password",
			args: args{
				ctx: context.Background(),
				userData: models.NewUser{
					Username: "abc",
					Email:    "abc@mail.com",
					Password: "abcd24",
				},
			},
			want: models.User{
				Username:     "abc",
				Email:        "abc@gmail.com",
				PasswordHash: "9980",
			},
			wantErr: false,
			mockRepoResponse: func() (models.User, error) {
				return models.User{
					Username:     "abc",
					Email:        "abc@gmail.com",
					PasswordHash: "9980",
				}, nil
			},
		},
		{
			name: "invalid hasedpass",
			args: args{
				ctx: context.Background(),
				userData: models.NewUser{
					Username: "abc",
					Email:    "abc@mail.com",
				},
			},
			want: models.User{
				// Username:     "abc",
				// Email:        "abc@gmail.com",
				// PasswordHash: "9980",
			},
			wantErr: true,
			mockRepoResponse: func() (models.User, error) {

				return models.User{}, errors.New("password hashing not done")

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, &auth.Auth{})
			got, err := s.UserSignup(tt.args.ctx, tt.args.userData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UserSignup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.UserSignup() = %v, want %v", got, tt.want)
			}
		})
	}
}
