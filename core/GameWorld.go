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

// constants for how big is the world
const (
	WorldWidth  = 800.0
	WorldHeight = 640.0
)

// GameObjectUpdate Type used for keeping track of new and deleted objects
type GameObjectUpdate struct {
	Type GameObjectType
	ID   uint64
}

// GameWorld the game world keeping track of all game objects and updating them
type GameWorld struct {
	PlayerShips []Spaceship
	PlayerAmmos []Projectile

	NewGameObjects     []GameObjectUpdate
	DeletedGameObjects []GameObjectUpdate

	nextID uint64

	playerShipData []playerShipData
}

// NewGameWorld Creates and intializes new game world
func NewGameWorld() GameWorld {
	result := GameWorld{}
	result.PlayerShips = make([]Spaceship, 0, 2)
	result.NewGameObjects = make([]GameObjectUpdate, 0)
	result.DeletedGameObjects = make([]GameObjectUpdate, 0)
	result.PlayerAmmos = make([]Projectile, 0, 20) // 20 sounds like reasonable number of ammos at the same time
	result.nextID = 0
	result.playerShipData = make([]playerShipData, 0, 2)

	return result
}

// AddNewPlayer adds new player to the world and returns the id of the player spaceship
func (world *GameWorld) AddNewPlayer() (id uint64) {
	id = world.getNextID()
	world.PlayerShips = append(world.PlayerShips, newPlayerShip(id, WorldWidth/2-PlayerShipWidth/2, WorldHeight-PlayerShipHeight))
	world.NewGameObjects = append(world.NewGameObjects, GameObjectUpdate{Type: PlayerShip, ID: id})

	world.playerShipData = append(world.playerShipData, playerShipData{})
	return
}

func (world *GameWorld) getNextID() (id uint64) {
	id = world.nextID
	world.nextID++
	return
}

// Update updates the game world
// it receives array of directions which the player spaceship took
// and updates every object positions
func (world *GameWorld) Update(dirs []SpaceshipDirection) {
	for i := range world.PlayerShips {
		select {
		case <-world.playerShipData[i].fire:
			obj := &world.PlayerShips[i].ObjectDimensions
			id := world.getNextID()
			world.PlayerAmmos = append(world.PlayerAmmos, newProjectile(id, obj.X+obj.width/2-AmmoWidth/2, obj.Y, PlayerAmmo, -10))
			world.NewGameObjects = append(world.NewGameObjects, GameObjectUpdate{Type: Ammo, ID: id})
		default:
		}
		movePlayerShip(&world.PlayerShips[i], dirs[i])
	}

	deleted := 0
	for i := range world.PlayerAmmos {
		j := i - deleted
		moveProjectile(&world.PlayerAmmos[j])
		// check for projectiles outside the world
		y := world.PlayerAmmos[j].Y
		id := world.PlayerAmmos[j].ID
		if y < -AmmoHeight {
			world.PlayerAmmos = world.PlayerAmmos[:j+copy(world.PlayerAmmos[j:], world.PlayerAmmos[j+1:])]
			world.DeletedGameObjects = append(world.DeletedGameObjects, GameObjectUpdate{Type: Ammo, ID: id})
			deleted++
		}
	}

	// check collisions

}

// Start starts the game
func (world *GameWorld) Start() {
	for i := range world.playerShipData {
		world.playerShipData[i].startSpawner()
	}
}

// Stop stops and cleans up the game
func (world *GameWorld) Stop() {
	for _, data := range world.playerShipData {
		data.cancelFire <- struct{}{}
	}
}

// ClearUpdates Clears arrays keeping track of new and deleted objects
func (world *GameWorld) ClearUpdates() {
	world.NewGameObjects = make([]GameObjectUpdate, 0)
	world.DeletedGameObjects = make([]GameObjectUpdate, 0)
}
