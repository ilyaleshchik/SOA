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
