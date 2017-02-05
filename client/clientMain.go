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

type entityCreateParams struct {
	texture string
	width   float32
	height  float32
}

// ship := gameEntity{BasicEntity: ecs.NewBasic()}
// 	ship.SpaceComponent = common.SpaceComponent{
// 		Position: engo.Point{X: 100, Y: 300},
// 		Width:    128,
// 		Height:   128,
// 	}
// 	texture, err := common.LoadedSprite("textures/ship.png")
// 	if err != nil {
// 		log.Println("Unable to load texture: " + err.Error())
// 	}

// 	ship.RenderComponent = common.RenderComponent{
// 		Drawable: texture,
// 		Scale:    engo.Point{X: ship.SpaceComponent.Width / texture.Width(), Y: ship.SpaceComponent.Height / texture.Height()},
// 	}

// 	for _, system := range world.Systems() {
// 		switch sys := system.(type) {
// 		case *common.RenderSystem:
// 			sys.Add(&ship.BasicEntity, &ship.RenderComponent, &ship.SpaceComponent)
// 		case *entityUpdater:
// 			sys.Add(&ship)
// 		}
// 	}

var entityCreateParamsMap map[core.EntityType]entityCreateParams

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

	common.SetBackground(color.White)
}

type gameEntity struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
}

type entityUpdater struct {
	entities map[uint64]*gameEntity
	decoder  *gob.Decoder
	world    *ecs.World
}

func (*entityUpdater) Remove(ecs.BasicEntity) {}
func (e *entityUpdater) Add(entity *gameEntity, id uint64) {
	e.entities[id] = entity
}

func (e *entityUpdater) Update(dt float32) {
	//fmt.Println("[Client] Receiving state")
	//var newPosition core.Spaceship
	var newWorld core.GameWorld
	e.decoder.Decode(&newWorld)

	// add new entities
	for _, en := range newWorld.NewEntities {
		params := entityCreateParamsMap[en.Type]
		entity := gameEntity{BasicEntity: ecs.NewBasic()}
		entity.SpaceComponent = common.SpaceComponent{
			Position: engo.Point{X: 0, Y: 0},
			Width:    params.width,
			Height:   params.height,
		}
		texture, err := common.LoadedSprite(params.texture)
		if err != nil {
			log.Println("Unable to load texture: " + err.Error())
		}

		entity.RenderComponent = common.RenderComponent{
			Drawable: texture,
			Scale:    engo.Point{X: entity.SpaceComponent.Width / texture.Width(), Y: entity.SpaceComponent.Height / texture.Height()},
		}

		for _, system := range e.world.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&entity.BasicEntity, &entity.RenderComponent, &entity.SpaceComponent)
			case *entityUpdater:
				sys.Add(&entity, en.ID)
			}
		}
	}

	// delete removed entitites
	// TODO: implement me

	// update entities
	for _, player := range newWorld.PlayerShips {
		space := &e.entities[player.ID].SpaceComponent
		space.Position.X = player.X
		space.Position.Y = player.Y
	}
}

func (e *entityUpdater) New(w *ecs.World) {
	e.entities = make(map[uint64]*gameEntity)
	e.world = w
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

func StartClient(serverIp string) {
	fmt.Println("[Client] starting with server ", serverIp)
	conn, err := net.Dial("tcp", serverIp+":1234")
	checkError(err)

	gameScene := gameScene{serverConnection: conn}

	// initialize entity create map
	entityCreateParamsMap = map[core.EntityType]entityCreateParams{
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
