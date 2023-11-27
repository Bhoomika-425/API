package repository

import (
	"context"
	"errors"
	"fmt"
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

// func (r *Repo) CheckEmail(ctx context.Context,email string)(models.Check,error){
//       var Checkmail models.Check
// 	  res:=r.DB.Where("email=?",email).First(&Checkmail)
// 	  if res.Error != nil{
// 		log.Info().Err(res.Error).Send()
// 		return models.Check{},errors.New("email is not valid")
// 	  }
// 	  return Checkmail,nil
// }

func (r *Repo) PasswordUpdating(ctx context.Context, details models.Valid, hashedpass string) error {
	var userdetails models.User

	result := r.DB.Where("email=?", details.Email).First(&userdetails)
	if result.Error != nil {
		log.Error().Err(result.Error).Str("email", details.Email).Msg("Error updating password")

	}

	userdetails.PasswordHash = hashedpass
	result = r.DB.Save(&userdetails)
	if result.Error != nil {
		// 	log.Info().Err(result.Error).Send()
		// 	return errors.New("password updation failed")
		// }
		log.Error().Err(result.Error).Str("email", details.Email).Msg("Error saving user details")
		return fmt.Errorf("failed to save user details: %w", result.Error)
	}
	return nil
}
