package types

import "net"

type Server struct {
	Port     int
	Listener net.Listener
}

type Connection struct {
	Id  int
	Con net.Conn
}
