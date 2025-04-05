package algorithms

import (
	"fmt"
	"slices"

	"github.com/fatih/color"
)

const PathNotFound = -1

// Node представляет собой узел графа
type Node struct {
	X, Y    int
	G, H, F int
	Parent  *Node
	VParent *Node // Visible parent
	Visited bool
	Index   int // Индекс для управления кучей
}

// IsValid проверяет, является ли позиция допустимой (на доске и не является стеной)
func IsValid(board [][]bool, x, y int) bool {
	return x >= 0 && x < len(board) && y >= 0 && y < len(board[0]) && !board[x][y]
}

// GetBoundaryCells возвращает массив всех пустых клеток на границе матрицы, отличных от стартовой клетки
func GetBoundaryCells(board [][]bool, startX, startY int) [][2]int {
	var boundaryCells [][2]int
	rows := len(board)
	cols := len(board[0])

	// Проверяем верхнюю границу
	for y := 0; y < cols; y++ {
		if IsValid(board, 0, y) && !(startX == 0 && startY == y) {
			boundaryCells = append(boundaryCells, [2]int{0, y})
		}
	}

	// Проверяем нижнюю границу
	for y := 0; y < cols; y++ {
		if IsValid(board, rows-1, y) && !(startX == rows-1 && startY == y) {
			boundaryCells = append(boundaryCells, [2]int{rows - 1, y})
		}
	}

	// Проверяем левую границу
	for x := 0; x < rows; x++ {
		if IsValid(board, x, 0) && !(startX == x && startY == 0) {
			boundaryCells = append(boundaryCells, [2]int{x, 0})
		}
	}

	// Проверяем правую границу
	for x := 0; x < rows; x++ {
		if IsValid(board, x, cols-1) && !(startX == x && startY == cols-1) {
			boundaryCells = append(boundaryCells, [2]int{x, cols - 1})
		}
	}

	return boundaryCells
}

// PrintBoard выводит матрицу лабиринта в терминал, выделяя путь
func PrintBoard(board [][]bool, path []Node, targets [][2]int, startX int, startY int) {
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
