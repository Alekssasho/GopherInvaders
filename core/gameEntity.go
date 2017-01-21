package core

import (
	"strconv"
)

// GameEntity is basic data for an enitty in the game
type GameEntity struct {
	ID   uint64
	X, Y float32
}

func (g *GameEntity) String() string {
	s := strconv.Itoa(int(g.ID))
	s += " "
	s += strconv.FormatFloat(float64(g.X), 'f', 2, 32)
	s += " "
	s += strconv.FormatFloat(float64(g.Y), 'f', 2, 32)

	return s
}
