package client

import (
	"encoding/gob"
	"image/color"
	"net"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type gameScene struct {
	serverConnection net.Conn
}

func (*gameScene) Type() string { return "Game Scene" }
func (*gameScene) Preload() {
	engo.Files.Load("textures/ship.png")
}

func (scene *gameScene) Setup(world *ecs.World) {
	engo.Input.RegisterButton("MoveLeft", engo.ArrowLeft)
	engo.Input.RegisterButton("MoveRight", engo.ArrowRight)
	engo.Input.RegisterButton("MoveUp", engo.ArrowUp)
	engo.Input.RegisterButton("MoveDown", engo.ArrowDown)

	world.AddSystem(&common.RenderSystem{})

	// setup the Updater
	updater := &entityUpdater{decoder: gob.NewDecoder(scene.serverConnection)}
	world.AddSystem(updater)

	inputController := &inputController{encoder: gob.NewEncoder(scene.serverConnection)}
	world.AddSystem(inputController)

	common.SetBackground(color.White)
}
