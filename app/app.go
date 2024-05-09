package app

import (
	"context"
	"errors"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/tgiv014/to/domains/config"
	"github.com/tgiv014/to/domains/link"
	"gorm.io/gorm"
	"tailscale.com/client/tailscale"
	"tailscale.com/tsnet"
)

type App struct {
	cfg config.Config

	lc    *tailscale.LocalClient
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

	if a.cfg.AuthKey == "" {
		return errors.New("auth key is required")
	}

	ts := tsnet.Server{
		AuthKey: a.cfg.AuthKey,
		Dir:     a.cfg.DataPath,
		Logf:    func(format string, args ...any) {},
	}
	defer ts.Close()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err = ts.Up(ctx)
	if err != nil {
		return err
	}

	ln, err := ts.Listen("tcp", ":80")
	if err != nil {
		return err
	}
	defer ln.Close()

	lc, err := ts.LocalClient()
	if err != nil {
		return err
	}
	a.lc = lc

	// Persistence
	a.db = db
	a.db.AutoMigrate(&link.Link{})

	// Services
	a.links = link.NewService(a.db)

	return a.Router().RunListener(ln)
}
