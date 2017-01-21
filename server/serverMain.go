package server

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"github.com/Alekssasho/GopherInvaders/core"
)

// Starts game server
func StartServer() {
	service := "127.0.0.1:1234"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		fmt.Println("server connection made")

		encoder := gob.NewEncoder(conn)
		decoder := gob.NewDecoder(conn)

		for n := 0; n < 10; n++ {
			var entity core.GameEntity
			fmt.Println("server receive data")
			err = decoder.Decode(&entity)
			checkError(err)
			fmt.Println(entity.String())
			fmt.Println("server send data")
			err = encoder.Encode(entity)
			checkError(err)
		}

		conn.Close()
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
