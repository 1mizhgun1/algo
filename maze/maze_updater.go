package maze

import (
	"fmt"
	"os"

	"algo/handlers/models"
	"github.com/pkg/errors"
)

func UpdateMaze(filename string, mazeMap [][]bool, points []models.Point) ([][]bool, error) {
	for _, point := range points {
		mazeMap[point.X][point.Y] = !mazeMap[point.X][point.Y]
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

func RestoreMaze(filename string, originalFilename string) ([][]bool, error) {
	mazeMap, err := ParseMaze(originalFilename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse original maze")
	}

	_, err = UpdateMaze(filename, mazeMap, []models.Point{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to update maze")
	}

	return mazeMap, nil
}
