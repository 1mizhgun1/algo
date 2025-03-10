package models

type SolveMazeInput struct {
	MazeID      int     `json:"labirint_id"`
	AlgorithmID int     `json:"algorithm_id"`
	Start       Point   `json:"start"`
	End         []Point `json:"end,omitempty"`
}

type SolveMazeOutput struct {
	Path []Tranzition `json:"path"`
	Time int          `json:"time"`
}

type Tranzition struct {
	Start Point `json:"start"`
	End   Point `json:"end"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}
