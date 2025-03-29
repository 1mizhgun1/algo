package maze

import (
	"fmt"
	"os"

	"algo/handlers/models"
	"github.com/pkg/errors"
)

func UpdateMaze(filename string, mazeMap [][]bool, points []models.Point) ([][]bool, error) {
	for _, point := range points {
		mazeMap[point.Y][point.X] = !mazeMap[point.Y][point.X]
	}

	if err := os.Remove(filename); err != nil {
		return mazeMap, errors.Wrap(err, "failed to remove maze")
	}

	newFile, err := os.Create(filename)
	if err != nil {
		return mazeMap, errors.Wrap(err, "failed to create file")
	}
	defer newFile.Close()

	for _, row := range mazeMap {
		for i, cell := range row {
			if cell {
				_, err = fmt.Fprint(newFile, "1")
			} else {
				_, err = fmt.Fprint(newFile, "0")
			}

			if i < len(row)-1 {
				_, err = fmt.Fprint(newFile, " ")
			}
		}
		_, err = fmt.Fprintln(newFile)
	}

	return mazeMap, errors.Wrap(err, "failed to write maze to file")
}
