package app

import (
	"github.com/charmbracelet/log"
	"github.com/glebarez/sqlite"
	"github.com/tgiv014/to/config"
	"github.com/tgiv014/to/domains/link"
	"gorm.io/gorm"
)

type App struct {
	cfg config.Config

	db *gorm.DB
}

func NewApp(cfg config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run() error {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	a.db = db

	a.db.AutoMigrate(&link.Link{})

	log.Info("Serving")
	return a.Router().Run()
}
