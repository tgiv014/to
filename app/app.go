package app

import (
	"github.com/glebarez/sqlite"
	"github.com/tgiv014/to/config"
	"github.com/tgiv014/to/domains/link"
	"gorm.io/gorm"
)

type App struct {
	cfg config.Config

	db    *gorm.DB
	links *link.Service
}

func NewApp(cfg config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run() error {
	db, err := gorm.Open(sqlite.Open(a.cfg.DBPath), &gorm.Config{})
	if err != nil {
		return err
	}

	// Persistence
	a.db = db
	a.db.AutoMigrate(&link.Link{})

	// Services
	a.links = link.NewService(a.db)

	return a.Router().Run()
}
