package database

import (
	"errors"
	"fmt"
	"project/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// db connection
func DbConnection() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate function will ONLY create tables, missing columns and missing indexes, and WON'T change existing column's type or delete unused columns
	err = db.Migrator().AutoMigrate(&models.User{})
	if err != nil {
		// If there is an error while migrating, log the error message and stop the program
		return nil, err
	}
	err = db.Migrator().AutoMigrate(&models.Company{})
	if err != nil {
		// If there is an error while migrating, log the error message and stop the program
		return nil, err
	}
	err = db.Migrator().AutoMigrate(&models.Jobs{}, &models.Location{}, &models.TechnologyStack{}, &models.WorkMode{}, &models.Qualification{}, &models.Shift{}, &models.Jobtype{})
	if err != nil {
		// If there is an error while migrating, log the error message and stop the program
		return nil, err
	}
	return db, nil

	// err = db.Migrator().AutoMigrate(&models.NewUserApplication{})
	// 	if err != nil {
	// 		// If there is an error while migrating, log the error message and stop the program
	// 		return nil, err
	// 	}
	// 	return db, nil

}

// passwoed hashing
func Passwordhashing(password string) (string, error) {

	if password == "" {
		return "", errors.New("error password must not be empty")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error in hashing the password : %w", err)
	}
	return string(hashedPass), nil

}

func HashedPassword(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
