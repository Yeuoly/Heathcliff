package server

import (
	"fmt"
	"net"
	"os"
	"strconv"

	types "github.com/Yeuoly/Heathcliff/types"
	"gopkg.in/yaml.v2"
)

const (
	SIGNAL_OVER = 0
)

var current_id = 0
var connection_port = 0
var cons chan types.Connection
var listeners []func(id int, msg []byte, n int) int

func init() {
	t := struct {
		Port int
	}{}

	conf, err := os.ReadFile("./conf/server.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(conf, &t)
	if err != nil {
		panic(err)
	}

	connection_port = t.Port
	cons = make(chan types.Connection, 5)
}

func worker(id int, con <-chan types.Connection) {
	for v := range con {
		buf := make([]byte, 8192)
		var status = true
		for status {
			n, err := v.Con.Read(buf)
			if err != nil {
				break
			}
			if n != 0 {
				for _, f := range listeners {
					if f(v.Id, buf, n) == SIGNAL_OVER {
						status = false
					}
				}
			}
		}

		v.Con.Close()
	}
}

func serve(server *types.Server) {
	for i := 0; i < 5; i++ {
		go worker(i, cons)
	}

	for {
		client, err := server.Listener.Accept()
		if err != nil {
			fmt.Printf("failed to accept connection]\n")
			continue
		}

		fmt.Printf("got connection from %s\n", client.RemoteAddr().String())

		connection := types.Connection{
			Id:  current_id,
			Con: client,
		}
		current_id++
		cons <- connection
	}
}

func Run() {
	v := types.Server{}
	v.Port = connection_port

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(v.Port))
	if err != nil {
		panic(err)
	}

	v.Listener = listener

	fmt.Printf("tcp server running on tcp:%d\n", v.Port)

	serve(&v)
}

func AppendListener(cb func(id int, msg []byte, n int) int) {
	listeners = append(listeners, cb)
}
