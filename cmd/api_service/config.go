package main

import (
	"log"
	"os"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

var config struct {
	Bind string `env:"BIND,required"`

	AccountsServiceHost string `env:"ACC_SERVICE_HOST,required"`
	PostsServiceHost    string `env:"POSTS_SERVICE_HOST,required"`

	DB      string `env:"DB,required"`
	DBDebug bool   `env:"DB_DEBUG,required"`

	PrivateSecret string `env:"PRIVATE_SECRET,required"`
	PublicSecret  string `env:"PUBLIC_SECRET,required"`
}

func loadConfig() {
	err := godotenv.Load(".env")
	if err != nil && !os.IsNotExist(err) {
		log.Panicf("[Config] .env file isn't exist %v", err)
	}

	if err = env.Parse(&config); err != nil {
		log.Panicf("[Config] Failed to parsing %v", err)
	}
}
