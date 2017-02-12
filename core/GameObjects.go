package core

// Spaceship is basic data for an enitty in the game
type Spaceship struct {
	ID   uint64
	X, Y float32
}

// SpaceshipDirection enum type for the direction of the player spaceship
type SpaceshipDirection int

// Enum types for different directions
const (
	Still SpaceshipDirection = iota
	Up
	Down
	Left
	Right
	UpRight
	UpLeft
	DownRight
	DownLeft
)
