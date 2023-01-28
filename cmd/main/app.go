package main

import (
	"exchanger/internal/core"
	"exchanger/internal/logger/zaplog"
	"flag"
	"github.com/joho/godotenv"
)

const (
	ENV = ".env"
)

func main() {
	var configFileName string
	flag.StringVar(&configFileName, "config", "dev", "config yml file path")
	flag.Parse()

	if err := godotenv.Load(ENV); err != nil {
		zaplog.AppLogger.Fatalf("[Error while load env]: %s", err.Error())
	}

	app := core.New(configFileName)
	if err := app.Serve(); err != nil {
		zaplog.AppLogger.Fatalf("[Error while core serv]: %s", err.Error())
	}
}
