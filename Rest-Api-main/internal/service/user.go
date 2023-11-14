package service

import (
	"context"
	"errors"
	"project/internal/database"
	"project/internal/models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func (s Service) UserLogin(ctx context.Context, userData models.NewUser) (string, error) {
	// checcking the email in the db
	var userDetails models.User
	userDetails, err := s.UserRepo.Userbyemail(ctx, userData.Email)
	if err != nil {
		return "", err
	}

	// comaparing the password and hashed password
	err = database.HashedPassword(userData.Password, userDetails.PasswordHash)
	if err != nil {
		log.Info().Err(err).Send()
		return "", errors.New("entered password is not wrong")
	}

	// setting up the claims
	claims := jwt.RegisteredClaims{
		Issuer:    "job portal project",
		Subject:   strconv.FormatUint(uint64(userDetails.ID), 10),
		Audience:  jwt.ClaimStrings{"users"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token, err := s.auth.GenerateToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (s Service) UserSignup(ctx context.Context, userData models.NewUser) (models.User, error) {
	hashedPass, err := database.Passwordhashing(userData.Password)
	if err != nil {
		return models.User{}, err
	}
	userDetails := models.User{
		Username:     userData.Username,
		Email:        userData.Email,
		PasswordHash: hashedPass,
	}
	userDetails, err = s.UserRepo.CreateUser(ctx, userDetails)
	if err != nil {
		return models.User{}, err
	}
	return userDetails, nil
}
