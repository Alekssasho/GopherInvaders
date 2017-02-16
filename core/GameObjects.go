package core

// ObjectDimensions is common struct for all game objects which needs position and size
type ObjectDimensions struct {
	ID   uint64
	X, Y float32

	width, height float32
}

// Spaceship is basic data for an enitty in the game
type Spaceship struct {
	ObjectDimensions
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

// ProjectileType enum type for ammo type
type ProjectileType int

// Different kind of ammos
const (
	PlayerAmmo = iota
	EnemyAmmo
)

// Projectile is ammo created by the ships
type Projectile struct {
	ObjectDimensions

	Type     ProjectileType
	velocity float32
	// TODO: add different types of ammo
}

// constants for how big are objects
const (
	PlayerShipWidth  = 64
	PlayerShipHeight = 64
	AmmoWidth        = 12
	AmmoHeight       = 12
)

func newPlayerShip(id uint64, x, y float32) Spaceship {
	return Spaceship{ObjectDimensions: ObjectDimensions{ID: id, X: x, Y: y, width: PlayerShipWidth, height: PlayerShipHeight}}
}

func newProjectile(id uint64, x, y float32, t ProjectileType, vel float32) Projectile {
	return Projectile{ObjectDimensions: ObjectDimensions{ID: id, X: x, Y: y, width: AmmoWidth, height: AmmoHeight}, Type: t, velocity: vel}
}

const (
	velocity         float32 = 5
	diagonalVelocity float32 = velocity / 1.4
)

func movePlayerShip(ship *Spaceship, dir SpaceshipDirection) {
	switch dir {
	case Up:
		ship.Y -= velocity
	case UpLeft:
		ship.X -= diagonalVelocity
		ship.Y -= diagonalVelocity
	case UpRight:
		ship.X += diagonalVelocity
		ship.Y -= diagonalVelocity
	case Down:
		ship.Y += velocity
	case DownLeft:
		ship.X -= diagonalVelocity
		ship.Y += diagonalVelocity
	case DownRight:
		ship.X += diagonalVelocity
		ship.Y += diagonalVelocity
	case Left:
		ship.X -= velocity
	case Right:
		ship.X += velocity
	}
}

func moveProjectile(proj *Projectile) {
	proj.Y += proj.velocity
}
