package main

import (
	"github.com/charmbracelet/log"
	"github.com/tgiv014/to/app"
	"github.com/tgiv014/to/domains/config"
)

func main() {
	log.Info("Loading config")
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	a := app.NewApp(cfg)
	err = a.Run()
	if err != nil {
		log.Fatal(err)
	}
}
