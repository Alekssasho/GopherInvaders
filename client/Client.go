package client

import (
	"fmt"
	"net"
	"os"

	"engo.io/engo"
	"github.com/Alekssasho/GopherInvaders/core"
)

// StartClient starts the game client. It receives ip of the server to connect to
func StartClient(serverIP string) {
	fmt.Println("[Client] starting with server ", serverIP)
	conn, err := net.Dial("tcp", serverIP+":1234")
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}

	gameScene := gameScene{serverConnection: conn}

	// initialize entity create map
	entityCreateParamsMap = map[core.GameObjectType]entityCreateParams{
		core.PlayerShip: entityCreateParams{texture: "textures/ship.png", width: 64, height: 64},
		// TODO: add other
	}

	opts := engo.RunOptions{
		Title:  "Gopher Invaders",
		Width:  800,
		Height: 640,
	}

	engo.Run(opts, &gameScene)
}
