package tailscale

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/tgiv014/to/domains/identity"
	tailscaleClient "tailscale.com/client/tailscale"
	"tailscale.com/tsnet"
)

type Service struct {
	authKey string
	dataDir string

	lc *tailscaleClient.LocalClient
}

func NewService(authKey, dataDir string) *Service {
	return &Service{
		authKey: authKey,
		dataDir: dataDir,
	}
}

func (s *Service) Listen() (net.Listener, error) {
	if s.authKey == "" {
		return nil, errors.New("auth key is required")
	}

	ts := tsnet.Server{
		AuthKey: s.authKey,
		Dir:     s.dataDir,
		Logf:    func(format string, args ...any) {},
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := ts.Up(ctx)
	if err != nil {
		return nil, err
	}

	s.lc, err = ts.LocalClient()
	if err != nil {
		return nil, err
	}

	ln, err := ts.Listen("tcp", ":80")
	if err != nil {
		return nil, err
	}

	return wrappedListener{ln, &ts}, nil
}

func (s *Service) Identify(c *gin.Context) (identity.User, error) {
	who, err := s.lc.WhoIs(c, c.Request.RemoteAddr)
	if err != nil {
		return identity.User{}, err
	}

	return identity.User{
		Name:          who.UserProfile.DisplayName,
		Login:         who.UserProfile.LoginName,
		ProfilePicURL: who.UserProfile.ProfilePicURL,
	}, nil
}

type wrappedListener struct {
	net.Listener
	tsServer *tsnet.Server
}

func (w wrappedListener) Close() error {
	err := w.Listener.Close()
	if err != nil {
		log.Error(err)
	}
	err = w.tsServer.Close()
	if err != nil {
		log.Error(err)
	}
	return nil
}
