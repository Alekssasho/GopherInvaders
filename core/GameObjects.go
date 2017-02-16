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

// Projectile is ammo created by the ships
type Projectile struct {
	ObjectDimensions

	Type     GameObjectType
	velocity float32
	// TODO: add different types of ammo
}

// constants for how big are objects
const (
	PlayerShipWidth  float32 = 64.0
	PlayerShipHeight float32 = 64.0
	AmmoWidth        float32 = 10.0
	AmmoHeight       float32 = 14.0
	EnemyWidth       float32 = 40.0
	EnemyHeight      float32 = 40.0
)

func newPlayerShip(id uint64, x, y float32) Spaceship {
	return Spaceship{ObjectDimensions: ObjectDimensions{ID: id, X: x, Y: y, width: PlayerShipWidth, height: PlayerShipHeight}}
}

func newProjectile(id uint64, x, y float32, t GameObjectType, vel float32) Projectile {
	return Projectile{ObjectDimensions: ObjectDimensions{ID: id, X: x, Y: y, width: AmmoWidth, height: AmmoHeight}, Type: t, velocity: vel}
}

func newEnemyShip(id uint64, x, y float32) Spaceship {
	return Spaceship{ObjectDimensions: ObjectDimensions{ID: id, X: x, Y: y, width: EnemyWidth, height: EnemyHeight}}
}

// This is per second
const (
	velocity         float32 = 250
	diagonalVelocity float32 = velocity / 1.4
)

func movePlayerShip(ship *Spaceship, dir SpaceshipDirection, delta float32) {
	switch dir {
	case Up:
		ship.Y -= velocity * delta
	case UpLeft:
		ship.X -= diagonalVelocity * delta
		ship.Y -= diagonalVelocity * delta
	case UpRight:
		ship.X += diagonalVelocity * delta
		ship.Y -= diagonalVelocity * delta
	case Down:
		ship.Y += velocity * delta
	case DownLeft:
		ship.X -= diagonalVelocity * delta
		ship.Y += diagonalVelocity * delta
	case DownRight:
		ship.X += diagonalVelocity * delta
		ship.Y += diagonalVelocity * delta
	case Left:
		ship.X -= velocity * delta
	case Right:
		ship.X += velocity * delta
	}
}

func moveProjectile(proj *Projectile, delta float32) {
	proj.Y += proj.velocity * delta
}
