package lazy_theta_star

import (
	"container/heap"
	"log"
	"math"
	"os"

	"algo/algorithms"
	"algo/maze"
	"github.com/pkg/errors"
)

// PriorityQueue реализует очередь приоритетов для узлов
type PriorityQueue []*algorithms.Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].F < pq[j].F
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*algorithms.Node)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // Устанавливаем nil для предотвращения утечки памяти
	item.Index = -1 // Устанавливаем индекс как -1 для предотвращения использования устаревших индексов
	*pq = old[0 : n-1]
	return item
}

// heuristic возвращает эвристическое значение между двумя узлами
func heuristic(a, b *algorithms.Node) int {
	return int(math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y)))
}

// isInClosedList проверяет, находится ли узел в закрытом списке
func isInClosedList(closedList map[[2]int]bool, node *algorithms.Node) bool {
	_, found := closedList[[2]int{node.X, node.Y}]
	return found
}

// isInOpenList проверяет, находится ли узел в открытом списке
func isInOpenList(openList map[[2]int]*algorithms.Node, node *algorithms.Node) bool {
	_, found := openList[[2]int{node.X, node.Y}]
	return found
}

// reconstructPath восстанавливает путь от целевого узла до стартового
func reconstructPath(current *algorithms.Node) []algorithms.Node {
	path := make([]algorithms.Node, 0)
	for current != nil {
		path = append([]algorithms.Node{*current}, path...)
		current = current.VParent
	}
	return path
}

// lineOfSight проверяет, есть ли прямая видимость между двумя узлами
func lineOfSight(board [][]bool, start, end *algorithms.Node) bool {
	if start == nil || end == nil {
		return false
	}

	x0, y0 := start.X, start.Y
	x1, y1 := end.X, end.Y
	dx, dy := x1-x0, y1-y0
	f := 0

	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}

	if dx >= dy {
		var y int
		if dy == 0 {
			y = y0
		} else {
			y = y0 + dy*(f+dx/2)/dx
		}
		for x := x0; x <= x1; x++ {
			if x != x0 || y != y0 {
				if !algorithms.IsValid(board, x, y) {
					return false
				}
			}
			f = f + dy
			if f >= dx {
				y = y + 1
				if (dy < 0 && y > y1) || (dy > 0 && y < y1) {
					return false
				}
				f = f - dx
			}
		}
	} else {
		var x int
		if dx == 0 {
			x = x0
		} else {
			x = x0 + dx*(f+dy/2)/dy
		}
		for y := y0; y <= y1; y++ {
			if x != x0 || y != y0 {
				if !algorithms.IsValid(board, x, y) {
					return false
				}
			}
			f = f + dx
			if f >= dy {
				x = x + 1
				if (dx < 0 && x > x1) || (dx > 0 && x < x1) {
					return false
				}
				f = f - dy
			}
		}
	}
	return true
}

// updateVertex обновляет значение узла
func updateVertex(openList *PriorityQueue, openListMap map[[2]int]*algorithms.Node, board [][]bool, node *algorithms.Node, neighbor *algorithms.Node) {
	if node.VParent == nil {
		return
	}

	if lineOfSight(board, node.VParent, neighbor) {
		newG := node.VParent.G + heuristic(node.VParent, neighbor)
		if newG < neighbor.G {
			neighbor.G = newG
			neighbor.F = newG + heuristic(neighbor, startNode)
			neighbor.VParent = node.VParent
			heap.Fix(openList, neighbor.Index)
		}
	} else {
		if neighbor.G > node.G+heuristic(node, neighbor) {
			neighbor.G = node.G + heuristic(node, neighbor)
			neighbor.F = neighbor.G + heuristic(neighbor, startNode)
			neighbor.VParent = node
			heap.Fix(openList, neighbor.Index)
		}
	}
}

var startNode *algorithms.Node

// LazyThetaStar алгоритм поиска кратчайшего пути
func LazyThetaStar(board [][]bool, startX, startY int, targets [][2]int) (int, []algorithms.Node) {
	openList := &PriorityQueue{}
	heap.Init(openList)
	closedList := make(map[[2]int]bool)
	startNode = &algorithms.Node{X: startX, Y: startY, G: 0, H: 0, F: 0, VParent: nil, Index: 0}
	heap.Push(openList, startNode)
	openListMap := make(map[[2]int]*algorithms.Node)
	openListMap[[2]int{startNode.X, startNode.Y}] = startNode

	for openList.Len() > 0 {
		current := heap.Pop(openList).(*algorithms.Node)
		delete(openListMap, [2]int{current.X, current.Y})
		closedList[[2]int{current.X, current.Y}] = true

		for _, target := range targets {
			if current.X == target[0] && current.Y == target[1] {
				return current.G, reconstructPath(current)
			}
		}

		neighbors := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
		for _, dir := range neighbors {
			neighbor := &algorithms.Node{X: current.X + dir[0], Y: current.Y + dir[1]}
			if !algorithms.IsValid(board, neighbor.X, neighbor.Y) {
				continue
			}

			if !isInClosedList(closedList, neighbor) {
				if !isInOpenList(openListMap, neighbor) {
					neighbor.G = math.MaxInt32
				}

				if neighbor.G > current.G+heuristic(current, neighbor) {
					neighbor.G = current.G + heuristic(current, neighbor)
					neighbor.F = neighbor.G + heuristic(neighbor, startNode)
					neighbor.VParent = current
					heap.Fix(openList, neighbor.Index)
				}

				if !isInOpenList(openListMap, neighbor) {
					heap.Push(openList, neighbor)
					openListMap[[2]int{neighbor.X, neighbor.Y}] = neighbor
				}

				updateVertex(openList, openListMap, board, current, neighbor)
			}
		}
	}

	return -1, nil
}

// TestLazyThetaStar тестирует алгоритм Lazy Theta*
func TestLazyThetaStar() {
	board, err := maze.ParseMaze(os.Getenv("MAZE_FILE_2"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to parse maze"))
	}

	startX, startY := 1, 0
	boundaryCells := algorithms.GetBoundaryCells(board, startX, startY)
	if len(boundaryCells) == 0 {
		log.Println("Нет доступных выходных клеток на границе")
		return
	}

	distance, path := LazyThetaStar(board, startX, startY, boundaryCells)
	if distance != -1 {
		log.Printf("Кратчайшее расстояние до выхода: %d\n", distance)
		log.Println("Лабиринт с выделенным путем:")
		algorithms.PrintBoard(board, path, boundaryCells, startX, startY)
	} else {
		log.Println("Путь до выхода не найден")
	}
}
