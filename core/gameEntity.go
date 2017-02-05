package core

// Spaceship is basic data for an enitty in the game
type Spaceship struct {
	ID   uint64
	X, Y float32
}

type SpaceshipDirection int

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

type EntityType int

const (
	PlayerShip EntityType = iota
	EnemyShip
	Ammo
)

type EntityUpdate struct {
	Type EntityType
	ID   uint64
}

type GameWorld struct {
	PlayerShips []Spaceship

	NewEntities     []EntityUpdate
	DeletedEntities []EntityUpdate

	nextId uint64
}

func NewGameWorld() GameWorld {
	result := GameWorld{}
	result.PlayerShips = make([]Spaceship, 0, 2)
	result.NewEntities = make([]EntityUpdate, 0)
	result.DeletedEntities = make([]EntityUpdate, 0)
	result.nextId = 0

	return result
}

func (world *GameWorld) AddNewPlayer() (id uint64) {
	id = world.nextId
	world.PlayerShips = append(world.PlayerShips, Spaceship{ID: id, X: 100, Y: 0})
	world.NewEntities = append(world.NewEntities, EntityUpdate{Type: PlayerShip, ID: id})
	world.nextId++
	return
}

func (world *GameWorld) Update(dirs []SpaceshipDirection) {
	const velocity float32 = 5
	const diagonalVelocity float32 = velocity / 1.4

	for i, _ := range world.PlayerShips {
		ship := &world.PlayerShips[i]
		switch dirs[i] {
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
}

func (world *GameWorld) ClearUpdates() {
	world.NewEntities = make([]EntityUpdate, 0)
	world.DeletedEntities = make([]EntityUpdate, 0)
}
