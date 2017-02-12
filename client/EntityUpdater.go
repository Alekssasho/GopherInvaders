package client

import (
	"encoding/gob"
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"github.com/Alekssasho/GopherInvaders/core"
)

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
	for _, en := range newWorld.NewGameObjects {
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
