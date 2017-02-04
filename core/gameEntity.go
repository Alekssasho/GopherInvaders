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

type GameWorld struct {
	PlayerShips []Spaceship
}

func NewGameWorld() GameWorld {
	result := GameWorld{}
	result.PlayerShips = make([]Spaceship, 0, 2)

	return result
}

func (world *GameWorld) AddNewPlayer() {
	world.PlayerShips = append(world.PlayerShips, Spaceship{ID: 0, X: 100, Y: 0})
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

// type List struct {
// 	Items []GameEntity
// }

// func (g *GameEntity) String() string {
// 	s := strconv.Itoa(int(g.ID))
// 	s += " "
// 	s += strconv.FormatFloat(float64(g.X), 'f', 2, 32)
// 	s += " "
// 	s += strconv.FormatFloat(float64(g.Y), 'f', 2, 32)

// 	return s
// }

// func (l *List) String() string {
// 	s := strconv.Itoa(len(l.Items))
// 	s += " { "
// 	for _, i := range l.Items {
// 		s += i.String()
// 		s += " "
// 	}
// 	s += "}"

// 	return s
// }
