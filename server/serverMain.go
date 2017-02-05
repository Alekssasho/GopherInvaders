package server

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"github.com/Alekssasho/GopherInvaders/core"
)

type player struct {
	encoder *gob.Encoder
	decoder *gob.Decoder
	id      uint64
}

// Starts game server
func StartServer(numPlayers int) {
	fmt.Println("[Server] starting")
	service := "0.0.0.0:1234"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	players := make([]player, 0, numPlayers)
	world := core.NewGameWorld()

	for i := 0; i < numPlayers; i++ {
		conn, err := listener.Accept()
		defer conn.Close()
		if err != nil {
			continue
		}

		fmt.Println("[Server] Connection made, adding new player")

		id := world.AddNewPlayer()

		players = append(players, player{
			encoder: gob.NewEncoder(conn),
			decoder: gob.NewDecoder(conn),
			id:      id})
	}

	fmt.Println("[Server] Starting game")
	for {
		// first we receive input from all client
		//fmt.Println("[Server] Receive dir")
		dirs := make([]core.SpaceshipDirection, numPlayers)
		for i, pl := range players {
			pl.decoder.Decode(&dirs[i])
		}

		// second update the world
		//fmt.Println("[Server] Update world")
		world.Update(dirs)

		// third send updated state to clients
		//fmt.Println("[Server] Send state")
		for _, pl := range players {
			pl.encoder.Encode(world)
		}

		world.ClearUpdates()
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
