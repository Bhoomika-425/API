package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"project/internal/auth"
	"project/internal/database"
	handler "project/internal/handlers"
	"project/internal/repository"
	service "project/internal/service"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

// Initilizaing the code
func main() {
	err := StartApp()
	if err != nil {
		log.Panic().Err(err).Send()
	}
	log.Info().Msg("this is our app:")
}

func StartApp() error {
	log.Info().Msg("intializing the authentication support")
	privatePEM, err := os.ReadFile(`C:\job-portal-api-main\private.pem`)
	if err != nil {
		return fmt.Errorf("error in reading auth private key : %w", err) // %w is used for error wraping
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return fmt.Errorf("error in parsing auth private key : %w", err) // %w is used for error wraping
	}
	publicPEM, err := os.ReadFile(`C:\job-portal-api-main\pubkey.pem`)
	if err != nil {
		return fmt.Errorf("error in reading auth public key : %w", err) // %w is used for error wraping
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return fmt.Errorf("error in parsing auth public key : %w", err) // %w is used for error wraping
	}
	a, err := auth.NewAuth(privateKey, publicKey)
	if err != nil {
		return fmt.Errorf("error in constructing auth %w", err)
	}
	log.Info().Msg("main started : initializing the data")

	db, err := database.DbConnection()
	if err != nil {
		return fmt.Errorf("error in opening the database connection : %w", err)
	}

	pg, err := db.DB()
	if err != nil {
		return fmt.Errorf("error in getting the database instance")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = pg.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("database is not connected: %w", err)
	}

	repo, err := repository.NewRepository(db)
	if err != nil {
		return err
	}

	sc, err := service.NewService(repo, a)
	if err != nil {
		return err
	}

	// initializing the http server
	api := http.Server{
		Addr:         ":8099",
		ReadTimeout:  8000 * time.Second,
		WriteTimeout: 800 * time.Second,
		IdleTimeout:  800 * time.Second,
		Handler:      handler.API(a, sc),
	}

	// channel to store any errors while setting up the service
	serverErrors := make(chan error, 1)

	go func() {
		log.Info().Str("Port", api.Addr).Msg("main started : api is listening")
		serverErrors <- api.ListenAndServe()
	}()

	//shutdown channel intercepts ctrl+c signals
	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, os.Interrupt)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error : %w", err)

	case sig := <-shutdown:
		log.Info().Msgf("main: Start shutdown %s", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			err := api.Close()
			return fmt.Errorf("could not stop server gracefully : %w", err)
		}
	}
	return nil

}
