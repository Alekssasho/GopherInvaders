package client

import (
	"encoding/gob"

	"engo.io/ecs"
	"engo.io/engo"
	"github.com/Alekssasho/GopherInvaders/core"
)

type inputController struct {
	encoder *gob.Encoder
}

func (i *inputController) Remove(basic ecs.BasicEntity) {
}

func (i *inputController) Update(dt float32) {
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
	i.encoder.Encode(dir)
}

func (i *inputController) Priority() int {
	return inputPriority
}
