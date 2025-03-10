package algorithms

const PathNotFound = -1

// Node представляет собой узел графа
type Node struct {
	X, Y    int
	G, H, F int
	Parent  *Node
}
