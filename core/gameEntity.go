package core

// Spaceship is basic data for an enitty in the game
type Spaceship struct {
	ID   uint64
	X, Y float32
}

type GameWorld struct {
	PlayerShips []Spaceship
}

func NewGameWorld() GameWorld {
	result := GameWorld{}
	result.PlayerShips = make([]Spaceship, 0, 2)

	return result
}

func (world *GameWorld) AddNewPlayer() {
	world.PlayerShips = append(world.PlayerShips, Spaceship{ID: 0, X: 100, Y: 200})
}

func (world *GameWorld) Update() {
	for i, _ := range world.PlayerShips {
		ship := &world.PlayerShips[i]
		ship.X += 5
		if ship.X >= 700 {
			ship.X = 10
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
