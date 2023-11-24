package config

import (
	"log"

	env "github.com/Netflix/go-env"
)

var cfg Config

type Config struct {
	AppConfig Appconfig
	DBconfig  DBconfig
	Redisconfig Redisconfig
	Authconfig Authconfig
}

type Appconfig struct {
	Port         string `env:"APP_PORT, required=true"`
	ReadTimeout  uint32 `env:"READ_TIMEOUT, required=true"`
	WriteTimeout uint32 `env:"WRITE_TIMEOUT, required=true"`
	IdleTimeout  uint32 `env:"IDLE_TIMEOUT, required=true"`
}

type DBconfig struct {
	// PostgresUser     string `"env:POSTGRES_USER, required=true"`
	// PostgresPassword string `"env:POSTGRES_PASSWORD, required=true"`
	// PostgresDb       string `"env:POSTGRES_DB, required=true"`
	DbCon string `env:"DB_DSN,required=true"`
}

type Redisconfig struct{
    Address string `env:"ADDRESS,required=true"`
	Password string `env:"PASS,required=true"`
	Database string `env:"DB,required=true"`
}

type Authconfig struct{
	PublicKey string `env:"PUBLICKEY,required=true"`
	PrivateKey string  `env:"PRIVATEKEY,required=true"`
}


func init() {
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		log.Panic(err)
	}
}
func GetConfig() Config {
	return cfg
}
