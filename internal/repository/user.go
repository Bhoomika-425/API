package repository

import (
	"context"
	"errors"
	"project/internal/models"

	"github.com/rs/zerolog/log"
)

func (r *Repo) CreateUser(ctx context.Context, UserDetails models.User) (models.User, error) {
	result := r.DB.Create(&UserDetails)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.User{}, errors.New("could not create the user")
	}
	return UserDetails, nil
}

func (r *Repo) Userbyemail(ctx context.Context, email string) (models.User, error) {
	var userDetails models.User
	result := r.DB.Where("email = ?", email).First(&userDetails)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.User{}, errors.New("email not found")
	}
	return userDetails, nil

}

func (r *Repo) CheckEmail(ctx context.Context,email string)(models.Check,error){
      var Checkmail models.Check
	  res:=r.DB.Where("email=?",email).First(&Checkmail)
	  if res.Error != nil{
		log.Info().Err(res.Error).Send()
		return models.Check{},errors.New("email is not valid")
	  }
	  return Checkmail,nil
}
