package goCommsNetDialer

import "net"

type iDialManager interface {
	Dial(network, addr string) (c net.Conn, err error)
}
