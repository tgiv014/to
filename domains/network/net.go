package network

import "net"

type Provider interface {
	Listen() (net.Listener, error)
}
