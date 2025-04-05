package a_star

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
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*algorithms.Node)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
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
		current = current.Parent
	}
	return path
}

// AStar алгоритм поиска кратчайшего пути
func AStar(board [][]bool, startX, startY int, targets [][2]int) (int, []algorithms.Node) {
	openList := &PriorityQueue{}
	heap.Init(openList)
	closedList := make(map[[2]int]bool)
	startNode := &algorithms.Node{X: startX, Y: startY, G: 0, H: 0, F: 0}
	heap.Push(openList, startNode)
	openListMap := make(map[[2]int]*algorithms.Node)
	openListMap[[2]int{startNode.X, startNode.Y}] = startNode

	for openList.Len() > 0 {
		current := heap.Pop(openList).(*algorithms.Node)
		delete(openListMap, [2]int{current.X, current.Y})

		for _, target := range targets {
			if current.X == target[0] && current.Y == target[1] {
				return current.G, reconstructPath(current)
			}
		}

		closedList[[2]int{current.X, current.Y}] = true

		neighbors := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
		for _, dir := range neighbors {
			neighbor := &algorithms.Node{X: current.X + dir[0], Y: current.Y + dir[1]}
			if !algorithms.IsValid(board, neighbor.X, neighbor.Y) || isInClosedList(closedList, neighbor) {
				continue
			}

			tentativeG := current.G + 1
			if !isInOpenList(openListMap, neighbor) {
				neighbor.H = heuristic(neighbor, &algorithms.Node{X: targets[0][0], Y: targets[0][1]})
				neighbor.F = tentativeG + neighbor.H
				neighbor.G = tentativeG
				neighbor.Parent = current
				heap.Push(openList, neighbor)
				openListMap[[2]int{neighbor.X, neighbor.Y}] = neighbor
			} else if tentativeG < neighbor.G {
				neighbor.G = tentativeG
				neighbor.F = neighbor.G + neighbor.H
				neighbor.Parent = current
				openListMap[[2]int{neighbor.X, neighbor.Y}] = neighbor
			}
		}
	}

	return algorithms.PathNotFound, nil
}

func TestAStar() {
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

	distance, path := AStar(board, startX, startY, boundaryCells)
	if distance != -1 {
		log.Printf("Кратчайшее расстояние до выхода: %d\n", distance)
		log.Println("Лабиринт с выделенным путем:")
		algorithms.PrintBoard(board, path, boundaryCells, startX, startY)
	} else {
		log.Println("Путь до выхода не найден")
	}
}
