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
	engo.Files.Load("textures/bullet.png")
	engo.Files.Load("textures/starfield.png")
	engo.Files.Load("textures/gopher.png")
	engo.Files.Load("fonts/DroidSerif.ttf")
}

func (scene *gameScene) Setup(world *ecs.World) {
	engo.Input.RegisterButton("MoveLeft", engo.ArrowLeft)
	engo.Input.RegisterButton("MoveRight", engo.ArrowRight)
	engo.Input.RegisterButton("MoveUp", engo.ArrowUp)
	engo.Input.RegisterButton("MoveDown", engo.ArrowDown)

	renderSystem := &common.RenderSystem{}
	world.AddSystem(renderSystem)

	// add font
	fnt := common.Font{
		URL:  "fonts/DroidSerif.ttf",
		FG:   color.White,
		Size: 18,
	}

	fnt.CreatePreloaded()

	fontUpdate := fontUpdater{BasicEntity: ecs.NewBasic(), font: fnt}
	fontUpdate.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 0, Y: windowHeight - 30},
	}
	fontUpdate.SetShader(common.HUDShader)

	// add background texture
	basicEntity := ecs.NewBasic()
	backgroundTexture, _ := common.LoadedSprite("textures/starfield.png")
	renderComponent := common.RenderComponent{
		Drawable: backgroundTexture,
		Scale:    engo.Point{X: windowWidth / backgroundTexture.Width(), Y: windowHeight / backgroundTexture.Height()},
	}
	spaceComponent := common.SpaceComponent{
		Position: engo.Point{X: 0, Y: 0},
		Width:    windowWidth,
		Height:   windowHeight,
	}
	renderSystem.Add(&basicEntity, &renderComponent, &spaceComponent)

	// setup the Updater
	updater := &entityUpdater{decoder: gob.NewDecoder(scene.serverConnection), fontUpdater: &fontUpdate}
	world.AddSystem(updater)

	inputController := &inputController{encoder: gob.NewEncoder(scene.serverConnection)}
	world.AddSystem(inputController)

	common.SetBackground(color.Black)

}
