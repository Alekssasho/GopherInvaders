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

// func StartClient() {
// 	// newEntity := core.GameEntity{BasicEntity: ecs.NewBasic()}
// 	// newEntity.SpaceComponent = common.SpaceComponent{
// 	// 	Position: engo.Point{10, 10},
// 	// 	Width:    64,
// 	// 	Height:   64,
// 	// }
// 	//newEntity := core.GameEntity{ID: 0, X: 10.0, Y: 52.0}
// 	newEntity := core.List{}
// 	newEntity.Items = make([]core.GameEntity, 0, 2)
// 	newEntity.Items = append(newEntity.Items, core.GameEntity{ID: 0, X: 10.0, Y: 52.0})
// 	newEntity.Items = append(newEntity.Items, core.GameEntity{ID: 1, X: 30.0, Y: 22.0})

// 	conn, err := net.Dial("tcp", "127.0.0.1:1234")
// 	checkError(err)

// 	fmt.Println("client connection made")

// 	encoder := gob.NewEncoder(conn)
// 	decoder := gob.NewDecoder(conn)

// 	for n := 0; n < 10; n++ {
// 		fmt.Println("client send data")
// 		err := encoder.Encode(newEntity)
// 		checkError(err)
// 		var entity core.List
// 		fmt.Println("client receive data")
// 		err = decoder.Decode(&entity)
// 		checkError(err)
// 		fmt.Println(entity.String())
// 	}

// 	os.Exit(0)

// }

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
	world.AddSystem(&common.RenderSystem{})

	// setup the Updater
	updater := &EntityUpdater{decoder: gob.NewDecoder(scene.serverConnection)}
	world.AddSystem(updater)

	// add the player spaceship
	ship := GameEntity{BasicEntity: ecs.NewBasic()}
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
		case *EntityUpdater:
			sys.Add(&ship)
		}
	}

	common.SetBackground(color.White)
}

type GameEntity struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
}

type EntityUpdater struct {
	entities []*GameEntity
	decoder  *gob.Decoder
}

func (*EntityUpdater) Remove(ecs.BasicEntity) {}
func (e *EntityUpdater) Add(entity *GameEntity) {
	e.entities = append(e.entities, entity)
}

func (e *EntityUpdater) Update(dt float32) {
	var newPosition core.Spaceship
	e.decoder.Decode(&newPosition)

	for _, entity := range e.entities {
		entity.SpaceComponent.Position.X = newPosition.X
		entity.SpaceComponent.Position.Y = newPosition.Y
	}
}

func (e *EntityUpdater) New(*ecs.World) {
	e.entities = make([]*GameEntity, 0, 10)
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
