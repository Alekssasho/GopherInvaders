package client

import (
	"engo.io/ecs"
	"engo.io/engo/common"
	"github.com/Alekssasho/GopherInvaders/core"
)

// Priorities for engo systems
const (
	updaterPriority = 10
	inputPriority   = 100 // Input needs to happen before update
)

type gameEntity struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
}

type entityCreateParams struct {
	texture string
	width   float32
	height  float32
}

// Map with creation parameters for different type of game objects
var entityCreateParamsMap map[core.GameObjectType]entityCreateParams
