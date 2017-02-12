package core

// GameObjectType is enum type for the type of the object
type GameObjectType int

const (
	// PlayerShip player controller ships
	PlayerShip GameObjectType = iota
	// EnemyShip enemy ships
	EnemyShip
	// Ammo weapon ammonutions
	Ammo
)

// GameObjectUpdate Type used for keeping track of new and deleted objects
type GameObjectUpdate struct {
	Type GameObjectType
	ID   uint64
}

// GameWorld the game world keeping track of all game objects and updating them
type GameWorld struct {
	PlayerShips []Spaceship

	NewGameObjects     []GameObjectUpdate
	DeletedGameObjects []GameObjectUpdate

	nextID uint64
}

// NewGameWorld Creates and intializes new game world
func NewGameWorld() GameWorld {
	result := GameWorld{}
	result.PlayerShips = make([]Spaceship, 0, 2)
	result.NewGameObjects = make([]GameObjectUpdate, 0)
	result.DeletedGameObjects = make([]GameObjectUpdate, 0)
	result.nextID = 0

	return result
}

// AddNewPlayer adds new player to the world and returns the id of the player spaceship
func (world *GameWorld) AddNewPlayer() (id uint64) {
	id = world.nextID
	world.PlayerShips = append(world.PlayerShips, Spaceship{ID: id, X: 100, Y: 0})
	world.NewGameObjects = append(world.NewGameObjects, GameObjectUpdate{Type: PlayerShip, ID: id})
	world.nextID++
	return
}

// Update updates the game world
// it receives array of directions which the player spaceship took
// and updates every object positions
func (world *GameWorld) Update(dirs []SpaceshipDirection) {
	const velocity float32 = 5
	const diagonalVelocity float32 = velocity / 1.4

	for i := range world.PlayerShips {
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

// ClearUpdates Clears arrays keeping track of new and deleted objects
func (world *GameWorld) ClearUpdates() {
	world.NewGameObjects = make([]GameObjectUpdate, 0)
	world.DeletedGameObjects = make([]GameObjectUpdate, 0)
}
