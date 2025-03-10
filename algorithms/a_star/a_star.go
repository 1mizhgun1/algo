package a_star

import (
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"slices"

	"algo/algorithms"
	"algo/maze"
	"github.com/fatih/color"
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

// isValid проверяет, является ли позиция допустимой (на доске и не является стеной)
func isValid(board [][]bool, x, y int) bool {
	return x >= 0 && x < len(board) && y >= 0 && y < len(board[0]) && !board[x][y]
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
			if !isValid(board, neighbor.X, neighbor.Y) || isInClosedList(closedList, neighbor) {
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

// GetBoundaryCells возвращает массив всех пустых клеток на границе матрицы, отличных от стартовой клетки
func GetBoundaryCells(board [][]bool, startX, startY int) [][2]int {
	var boundaryCells [][2]int
	rows := len(board)
	cols := len(board[0])

	// Проверяем верхнюю границу
	for y := 0; y < cols; y++ {
		if isValid(board, 0, y) && !(startX == 0 && startY == y) {
			boundaryCells = append(boundaryCells, [2]int{0, y})
		}
	}

	// Проверяем нижнюю границу
	for y := 0; y < cols; y++ {
		if isValid(board, rows-1, y) && !(startX == rows-1 && startY == y) {
			boundaryCells = append(boundaryCells, [2]int{rows - 1, y})
		}
	}

	// Проверяем левую границу
	for x := 0; x < rows; x++ {
		if isValid(board, x, 0) && !(startX == x && startY == 0) {
			boundaryCells = append(boundaryCells, [2]int{x, 0})
		}
	}

	// Проверяем правую границу
	for x := 0; x < rows; x++ {
		if isValid(board, x, cols-1) && !(startX == x && startY == cols-1) {
			boundaryCells = append(boundaryCells, [2]int{x, cols - 1})
		}
	}

	return boundaryCells
}

// printBoard выводит матрицу лабиринта в терминал, выделяя путь
func printBoard(board [][]bool, path []algorithms.Node, targets [][2]int, startX int, startY int) {
	pathMap := make(map[[2]int]bool)
	for _, node := range path {
		pathMap[[2]int{node.X, node.Y}] = true
	}

	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	for i, row := range board {
		for j, cell := range row {
			if slices.Contains(targets, [2]int{i, j}) {
				fmt.Printf("%s ", green("0"))
			} else if i == startX && j == startY {
				fmt.Printf("%s ", blue("0"))
			} else if pathMap[[2]int{i, j}] {
				fmt.Printf("%s ", red("0"))
			} else if cell {
				fmt.Print("1 ")
			} else {
				fmt.Print("0 ")
			}
		}
		fmt.Println()
	}
}

func TestAStar() {
	board, err := maze.ParseMaze(os.Getenv("MAZE_FILE_2"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to parse maze"))
	}

	startX, startY := 1, 0

	boundaryCells := GetBoundaryCells(board, startX, startY)
	if len(boundaryCells) == 0 {
		log.Println("Нет доступных выходных клеток на границе")
		return
	}

	distance, path := AStar(board, startX, startY, boundaryCells)
	if distance != -1 {
		log.Printf("Кратчайшее расстояние до выхода: %d\n", distance)
		log.Println("Лабиринт с выделенным путем:")
		printBoard(board, path, boundaryCells, startX, startY)
	} else {
		log.Println("Путь до выхода не найден")
	}
}
