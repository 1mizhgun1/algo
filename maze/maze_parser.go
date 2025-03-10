package maze

import (
	"bufio"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func ParseMaze(filename string) ([][]bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open maze file")
	}
	defer file.Close()

	var maze [][]bool

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var row []bool
		for _, char := range strings.Fields(line) {
			row = append(row, char == "1")
		}
		maze = append(maze, row)
	}

	if err = scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to scan maze file")
	}

	return maze, nil
}
