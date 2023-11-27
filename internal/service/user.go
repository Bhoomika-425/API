package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
	"project/internal/database"
	"project/internal/models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func (s Service) UserLogin(ctx context.Context, userData models.UserLogin) (string, error) {
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
		Dob:          userData.Dob,
	}
	userDetails, err = s.UserRepo.CreateUser(ctx, userDetails)
	if err != nil {
		return models.User{}, err
	}
	return userDetails, nil
}

func (s Service) ValidatingEmail(ctx context.Context, useremail models.Check) (string, error) {
	_, err := s.UserRepo.Userbyemail(ctx, useremail.Email)
	if err != nil {
		return "", errors.New("email not matched")
	}
	a := GenerateOTP(useremail.Email)
	return a, nil

}

func GenerateOTP(email string) string {
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	otp := rand.Intn(max-min+1) + min
	verificationcode := strconv.Itoa(otp)
	from := "bhoomikasabalur@gmail.com"
	password := "hghf uxjo oqgf qmzx"

	// Recipient's email address
	to := email

	// SMTP server details
	smtpServer := "smtp.gmail.com"
	smtpPort := 587

	// Message content
	message := []byte(fmt.Sprintf("Subject: Test Email\n\nThis is a test email body.", verificationcode))

	// Authentication information
	auth := smtp.PlainAuth("", from, password, smtpServer)

	// SMTP connection
	smtpAddr := fmt.Sprintf("%s:%d", smtpServer, smtpPort)
	err := smtp.SendMail(smtpAddr, auth, from, []string{to}, message)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return ""
	}

	fmt.Println("Email sent successfully!")

	return "successfully sent otp"

}

func (s Service) RecoveringPassword(ctx context.Context, Recoverpass models.Valid) error {
	verifycode, err := s.rdb.Getaotp(ctx, Recoverpass.Email)
	if err != nil {
		log.Error().Err(err).Msg("error retrieving otp from redis")
		return err
	}

	if verifycode != Recoverpass.Verifiedotp {
		return errors.New("Verification of otp is not correct")
	}

	if Recoverpass.Password != Recoverpass.Confirmpass {
		return errors.New("password &confirmpass is not matching")
	}

	str, err := database.Passwordhashing(Recoverpass.Confirmpass)
	if err != nil {
		log.Error().Err(err).Msg("password hashing is not done properly")
		return err
	}

	err = s.UserRepo.PasswordUpdating(ctx, Recoverpass, str)
	if err != nil {
		log.Error().Err(err).Msg("error in updating password to repo")
		return err
	}
	return nil

}
