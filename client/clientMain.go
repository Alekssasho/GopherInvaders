package client

import (
	"fmt"
	"log"
	"net"
	"os"

	"image/color"

	"encoding/gob"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"github.com/Alekssasho/GopherInvaders/core"
)

const (
	updaterPriority = 10
	inputPriority   = 100 // Input needs to happen begore update
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

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
	// engo.Input.RegisterButton("Fire", engo.Space)

	world.AddSystem(&common.RenderSystem{})

	// setup the Updater
	updater := &entityUpdater{decoder: gob.NewDecoder(scene.serverConnection)}
	world.AddSystem(updater)

	inputController := &inputController{encoder: gob.NewEncoder(scene.serverConnection)}
	world.AddSystem(inputController)

	// add the player spaceship
	ship := gameEntity{BasicEntity: ecs.NewBasic()}
	ship.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 100, Y: 300},
		Width:    128,
		Height:   128,
	}
	texture, err := common.LoadedSprite("textures/ship.png")
	if err != nil {
		log.Println("Unable to load texture: " + err.Error())
	}

	ship.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{X: ship.SpaceComponent.Width / texture.Width(), Y: ship.SpaceComponent.Height / texture.Height()},
	}

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&ship.BasicEntity, &ship.RenderComponent, &ship.SpaceComponent)
		case *entityUpdater:
			sys.Add(&ship)
		}
	}

	common.SetBackground(color.White)
}

type gameEntity struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
}

type entityUpdater struct {
	entities []*gameEntity
	decoder  *gob.Decoder
}

func (*entityUpdater) Remove(ecs.BasicEntity) {}
func (e *entityUpdater) Add(entity *gameEntity) {
	e.entities = append(e.entities, entity)
}

func (e *entityUpdater) Update(dt float32) {
	//fmt.Println("[Client] Receiving state")
	var newPosition core.Spaceship
	e.decoder.Decode(&newPosition)

	for _, entity := range e.entities {
		entity.SpaceComponent.Position.X = newPosition.X
		entity.SpaceComponent.Position.Y = newPosition.Y
	}
}

func (e *entityUpdater) New(*ecs.World) {
	e.entities = make([]*gameEntity, 0, 10)
}

func (e *entityUpdater) Priority() int {
	return updaterPriority
}

type inputController struct {
	entity  *gameEntity
	encoder *gob.Encoder
}

func (*inputController) Remove(ecs.BasicEntity) {}
func (e *inputController) Add(entity *gameEntity) {
	e.entity = entity
}

func (e *inputController) Update(dt float32) {
	dir := core.Still

	if engo.Input.Button("MoveLeft").Down() {
		dir = core.Left
	} else if engo.Input.Button("MoveRight").Down() {
		dir = core.Right
	}

	if engo.Input.Button("MoveUp").Down() {
		if dir == core.Left {
			dir = core.UpLeft
		} else if dir == core.Right {
			dir = core.UpRight
		} else {
			dir = core.Up
		}
	} else if engo.Input.Button("MoveDown").Down() {
		if dir == core.Left {
			dir = core.DownLeft
		} else if dir == core.Right {
			dir = core.DownRight
		} else {
			dir = core.Down
		}
	}

	//fmt.Println("[Client] Sending direction", dir)
	e.encoder.Encode(dir)
}

func (e *inputController) Priority() int {
	return inputPriority
}

func StartClient() {
	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	checkError(err)

	gameScene := gameScene{serverConnection: conn}

	opts := engo.RunOptions{
		Title:  "Gopher Invaders",
		Width:  800,
		Height: 640,
	}

	engo.Run(opts, &gameScene)
}
