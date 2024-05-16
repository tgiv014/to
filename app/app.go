package app

import (
	"github.com/glebarez/sqlite"
	"github.com/tgiv014/to/components"
	"github.com/tgiv014/to/domains/config"
	"github.com/tgiv014/to/domains/identity"
	"github.com/tgiv014/to/domains/link"
	"github.com/tgiv014/to/domains/network"
	"github.com/tgiv014/to/domains/tailscale"
	"gorm.io/gorm"
)

type App struct {
	Cfg config.Config

	Net        network.Provider
	Identifier identity.Identifier
	DB         *gorm.DB
	Links      *link.Service
}

func NewApp(cfg config.Config) *App {
	a := &App{Cfg: cfg}
	if a.Cfg.LocalPort != 0 {
		provider := network.NewLocalProvider(a.Cfg.LocalPort)
		a.Net = provider
		a.Identifier = provider
	} else {
		ts := tailscale.NewService(a.Cfg.AuthKey, a.Cfg.DataPath)
		a.Net = ts
		a.Identifier = ts
	}

	if a.Cfg.LiveReload {
		components.EnableLiveReload = true
	}
	return a
}

func (a *App) Setup() error {
	db, err := gorm.Open(sqlite.Open(a.Cfg.DBPath), &gorm.Config{})
	if err != nil {
		return err
	}
	a.DB = db

	err = a.DB.AutoMigrate(&link.Link{})
	if err != nil {
		return err
	}

	// Services
	a.Links = link.NewService(a.DB)

	return nil
}

func (a *App) Run() error {
	ln, err := a.Net.Listen()
	if err != nil {
		return err
	}
	return a.Router().RunListener(ln)
}
