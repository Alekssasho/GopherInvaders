package core

import (
	"math"
	"math/rand"
)

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

	EnemyShips []Spaceship

	NewGameObjects     []GameObjectUpdate
	DeletedGameObjects []GameObjectUpdate

	PlayerScore int

	nextID uint64

	playerShipData []playerShipData
	enemyShipData  []enemyShipData
	currentTime    float32
}

// NewGameWorld Creates and intializes new game world
func NewGameWorld() GameWorld {
	result := GameWorld{}
	result.PlayerShips = make([]Spaceship, 0, 2)
	result.NewGameObjects = make([]GameObjectUpdate, 0)
	result.DeletedGameObjects = make([]GameObjectUpdate, 0)
	result.PlayerAmmos = make([]Projectile, 0, 20) // 20 sounds like reasonable number of ammos at the same time
	result.EnemyShips = make([]Spaceship, 0, 20)
	result.PlayerScore = 0
	result.nextID = 0
	result.playerShipData = make([]playerShipData, 0, 2)
	result.enemyShipData = make([]enemyShipData, 0, 20)
	result.currentTime = 0.0

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
func (world *GameWorld) Update(dirs []SpaceshipDirection, delta float32) {
	for i := range world.PlayerShips {
		select {
		case <-world.playerShipData[i].fire:
			obj := &world.PlayerShips[i].ObjectDimensions
			id := world.getNextID()
			world.PlayerAmmos = append(world.PlayerAmmos, newProjectile(id, obj.X+obj.width/2-AmmoWidth/2, obj.Y, PlayerAmmo, -600))
			world.NewGameObjects = append(world.NewGameObjects, GameObjectUpdate{Type: Ammo, ID: id})
		default:
		}
		movePlayerShip(&world.PlayerShips[i], dirs[i], delta)
	}

	if len(world.EnemyShips) == 0 {
		// Time to add new enemies
		// Summon pattern is 12 enemies per row
		numRows := 2
		if world.PlayerScore > 0 {
			numRows += int(math.Ceil(math.Log10(float64(world.PlayerScore))))
		}
		spaceBetweenShip := (1.5 * EnemyWidth) / 5.0
		for i := 0; i < numRows; i++ {
			currentY := -EnemyHeight - float32(i)*(EnemyHeight*1.2)
			currentX := 2 * EnemyWidth
			rowDelta := rand.Float32()
			for j := 0; j < 12; j++ {
				if j == 6 {
					currentX -= spaceBetweenShip
					currentX += 1 * EnemyWidth
				}

				enemy := newEnemyShip(world.getNextID(), currentX, currentY)
				world.EnemyShips = append(world.EnemyShips, enemy)
				world.enemyShipData = append(world.enemyShipData, enemyShipData{startX: enemy.X, delta: rowDelta})
				world.NewGameObjects = append(world.NewGameObjects, GameObjectUpdate{Type: EnemyShip, ID: enemy.ID})
				currentX += spaceBetweenShip + EnemyWidth
			}
		}
		world.currentTime = 0.0
	} else {
		// enemy moves like a sin function
		world.currentTime += delta
		for i := range world.EnemyShips {
			obj := &world.EnemyShips[i]
			obj.Y += 15.0 * delta
			obj.X = world.enemyShipData[i].startX + float32(math.Sin(float64(world.currentTime*0.5+world.enemyShipData[i].delta)))*2*EnemyWidth
		}
	}

	deleted := 0
	firstEnemyY := world.EnemyShips[0].Y + EnemyHeight // used for early out
	for i := range world.PlayerAmmos {
		j := i - deleted
		moveProjectile(&world.PlayerAmmos[j], delta)
		// check for projectiles outside the world
		y := world.PlayerAmmos[j].Y
		id := world.PlayerAmmos[j].ID
		destroy := y < -AmmoHeight
		if !destroy {
			// check collision with objects
			if y > firstEnemyY {
				continue
			}

			for enemyIndex := range world.EnemyShips {
				if checkCollision(&world.PlayerAmmos[j].ObjectDimensions, &world.EnemyShips[enemyIndex].ObjectDimensions) {
					destroy = true
					world.PlayerScore++
					world.DeletedGameObjects = append(world.DeletedGameObjects, GameObjectUpdate{Type: EnemyShip, ID: world.EnemyShips[enemyIndex].ID})
					world.EnemyShips = world.EnemyShips[:enemyIndex+copy(world.EnemyShips[enemyIndex:], world.EnemyShips[enemyIndex+1:])]
					world.enemyShipData = world.enemyShipData[:enemyIndex+copy(world.enemyShipData[enemyIndex:], world.enemyShipData[enemyIndex+1:])]
					break
				}
			}
		}

		if destroy {
			world.PlayerAmmos = world.PlayerAmmos[:j+copy(world.PlayerAmmos[j:], world.PlayerAmmos[j+1:])]
			world.DeletedGameObjects = append(world.DeletedGameObjects, GameObjectUpdate{Type: Ammo, ID: id})
			deleted++
		}
	}
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
