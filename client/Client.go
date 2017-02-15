package client

import (
	"fmt"
	"net"
	"os"

	"engo.io/engo"
	"github.com/Alekssasho/GopherInvaders/core"
)

const (
	windowWidth  = 800.0
	windowHeight = 640.0
)

const (
	worldScaleX = windowWidth / core.WorldWidth
	worldScaleY = windowHeight / core.WorldHeight
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
		core.PlayerShip: entityCreateParams{texture: "textures/ship.png", width: core.PlayerShipWidth * worldScaleX, height: core.PlayerShipHeight * worldScaleY},
		core.Ammo:       entityCreateParams{texture: "textures/ship.png", width: core.AmmoWidth * worldScaleX, height: core.AmmoHeight * worldScaleY},
		// TODO: add other
	}

	opts := engo.RunOptions{
		Title:  "Gopher Invaders",
		Width:  windowWidth,
		Height: windowHeight,
	}

	engo.Run(opts, &gameScene)
}
