package models

import (
	"fmt"
	"time"

	"algo/config"
	"github.com/pkg/errors"
)

type SolveMazeInput struct {
	MazeID      int     `json:"labirint_id"`
	AlgorithmID int     `json:"algorithm_id"`
	Start       Point   `json:"start"`
	End         []Point `json:"end,omitempty"`
}

type SolveMazeOutput struct {
	Path          []Tranzition  `json:"path"`
	Dist          int           `json:"dist"`
	ExecutionTime time.Duration `json:"time"`
}

type Tranzition struct {
	Start Point `json:"start"`
	End   Point `json:"end"`
}

type Point struct {
	Y int `json:"x"`
	X int `json:"y"`
}

type UpdateMazeInput struct {
	MazeID int     `json:"labirint_id"`
	Points []Point `json:"points"`
}

type UpdateMazeOutput struct {
	Map [][]int `json:"labirint"`
}

type GetMazeInput struct {
	MazeID int `json:"labirint_id"`
}

type GetMazeOutput struct {
	Map [][]int `json:"labirint"`
}

type RestoreMazeInput struct {
	MazeID int `json:"labirint_id"`
}

type RestoreMazeOutput struct {
	Map [][]int `json:"labirint"`
}

func validateMazeID(mazeID int, cfg config.AppConfig) bool {
	return 0 < mazeID && mazeID <= cfg.MazeCount
}

func validateAlgorithmID(algorithmID int, cfg config.AppConfig) bool {
	return 0 < algorithmID && algorithmID <= cfg.AlgorithmCount
}

func validatePoint(point Point, n int, m int) bool {
	return 0 <= point.X && point.X < m && 0 <= point.Y && point.Y < n
}

func (req *SolveMazeInput) Validate(cfg config.AppConfig, n int, m int) error {
	if !validateMazeID(req.MazeID, cfg) {
		return errors.New("invalid labirint_id")
	}

	if !validateAlgorithmID(req.AlgorithmID, cfg) {
		return errors.New("invalid algorithm_id")
	}

	if !validatePoint(req.Start, n, m) {
		return errors.New("invalid start point")
	}

	for i, end := range req.End {
		if !validatePoint(end, n, m) {
			return fmt.Errorf("invalid end point at index %d", i)
		}
	}

	return nil
}

func (req *UpdateMazeInput) Validate(cfg config.AppConfig, n int, m int) error {
	if !validateMazeID(req.MazeID, cfg) {
		return errors.New("invalid labirint_id")
	}

	for i, point := range req.Points {
		if !validatePoint(point, n, m) {
			return fmt.Errorf("invalid point at index %d", i)
		}
	}

	return nil
}

func (req *GetMazeInput) Validate(cfg config.AppConfig) error {
	if !validateMazeID(req.MazeID, cfg) {
		return errors.New("invalid labirint_id")
	}

	return nil
}

func (req *RestoreMazeInput) Validate(cfg config.AppConfig) error {
	if !validateMazeID(req.MazeID, cfg) {
		return errors.New("invalid labirint_id")
	}

	return nil
}
