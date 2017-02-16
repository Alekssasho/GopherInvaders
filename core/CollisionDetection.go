package core

func checkCollision(lhs, rhs *ObjectDimensions) bool {
	return lhs.X < rhs.X+rhs.width &&
		lhs.X+lhs.width > rhs.X &&
		lhs.Y < rhs.Y+rhs.height &&
		lhs.height+lhs.Y > rhs.Y
}
