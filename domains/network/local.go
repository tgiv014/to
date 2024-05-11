package network

import (
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/tgiv014/to/domains/identity"
)

type LocalProvider struct {
	port int
}

func NewLocalProvider(port int) *LocalProvider {
	return &LocalProvider{
		port: port,
	}
}

func (l *LocalProvider) Listen() (net.Listener, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", l.port))
	if err != nil {
		return nil, err
	}

	return ln, nil
}

func (l *LocalProvider) Identify(c *gin.Context) (identity.User, error) {
	return identity.User{
		Name:  "Local User",
		Login: "local",
	}, nil
}
