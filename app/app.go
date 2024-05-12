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
	cfg config.Config

	netProvider network.Provider
	identifier  identity.Identifier
	db          *gorm.DB
	links       *link.Service
}

func NewApp(cfg config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run() error {
	if a.cfg.LocalPort != 0 {
		provider := network.NewLocalProvider(a.cfg.LocalPort)
		a.netProvider = provider
		a.identifier = provider
	} else {
		ts := tailscale.NewService(a.cfg.AuthKey, a.cfg.DataPath)
		a.netProvider = ts
		a.identifier = ts
	}

	if a.cfg.LiveReload {
		components.EnableLiveReload = true
	}

	ln, err := a.netProvider.Listen()
	if err != nil {
		return err
	}

	db, err := gorm.Open(sqlite.Open(a.cfg.DBPath), &gorm.Config{})
	if err != nil {
		return err
	}

	// Persistence
	a.db = db
	a.db.AutoMigrate(&link.Link{})

	// Services
	a.links = link.NewService(a.db)

	return a.Router().RunListener(ln)
}
