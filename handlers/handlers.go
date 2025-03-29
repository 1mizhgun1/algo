package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"algo/algorithms"
	"algo/algorithms/a_star"
	"algo/config"
	"algo/handlers/models"
	"algo/maze"
	"algo/utils"
)

type App struct {
	cfg config.AppConfig
}

func NewApp(cfg config.AppConfig) *App {
	return &App{cfg}
}

func (app *App) SolveMazeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.SolveMazeInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.LogError(ctx, err, utils.MsgErrUnmarshalRequest)
		http.Error(w, utils.Invalid, http.StatusBadRequest)
		return
	}

	board, err := maze.ParseMaze(os.Getenv(fmt.Sprintf("MAZE_FILE_%d", req.MazeID)))
	if err != nil {
		utils.LogError(ctx, err, "failed to parse maze")
		http.Error(w, utils.Internal, http.StatusInternalServerError)
		return
	}

	if err = req.Validate(app.cfg, len(board[0]), len(board)); err != nil {
		utils.LogError(ctx, err, "failed to validate maze")
		http.Error(w, utils.Invalid, http.StatusBadRequest)
		return
	}

	if board[req.Start.X][req.Start.Y] {
		utils.LogErrorMessage(ctx, fmt.Sprintf("start point (%d,%d)=1, it is wall", req.Start.X, req.Start.Y))
		http.Error(w, utils.Invalid, http.StatusBadRequest)
		return
	}

	boundaryCells := make([][2]int, 0)
	if len(req.End) == 0 {
		boundaryCells = a_star.GetBoundaryCells(board, req.Start.X, req.Start.Y)
	} else {
		for _, end := range req.End {
			if board[end.X][end.Y] {
				utils.LogErrorMessage(ctx, fmt.Sprintf("end point (%d,%d)=1, it is wall", end.X, end.Y))
				http.Error(w, utils.Invalid, http.StatusBadRequest)
				return
			}
			boundaryCells = append(boundaryCells, [2]int{end.X, end.Y})
		}
	}

	var (
		result   models.SolveMazeOutput
		distance int
		path     []algorithms.Node
	)

	switch req.AlgorithmID {
	case 1:
		distance, path = a_star.AStar(board, req.Start.X, req.Start.Y, boundaryCells)
	default:
		utils.LogErrorMessage(ctx, fmt.Sprintf("invalid algorithm id=%d", req.AlgorithmID))
		http.Error(w, utils.Invalid, http.StatusBadRequest)
		return
	}

	if distance == algorithms.PathNotFound {
		if err = json.NewEncoder(w).Encode(result); err != nil {
			utils.LogError(ctx, err, utils.MsgErrMarshalResponse)
			http.Error(w, utils.Internal, http.StatusInternalServerError)
			return
		}
		return
	}

	if len(path) == 0 {
		utils.LogErrorMessage(ctx, "path is empty")
		http.Error(w, utils.Internal, http.StatusInternalServerError)
		return
	}

	result.Time = distance
	result.Path = make([]models.Tranzition, len(path)-1)

	for i := 1; i < len(path); i++ {
		result.Path[i-1] = models.Tranzition{
			Start: models.Point{X: path[i-1].X, Y: path[i-1].Y},
			End:   models.Point{X: path[i].X, Y: path[i].Y},
		}
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		utils.LogError(ctx, err, utils.MsgErrMarshalResponse)
		http.Error(w, utils.Internal, http.StatusInternalServerError)
		return
	}
}

func (app *App) UpdateMazeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.UpdateMazeInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.LogError(ctx, err, utils.MsgErrUnmarshalRequest)
		http.Error(w, utils.Invalid, http.StatusBadRequest)
		return
	}

	filename := os.Getenv(fmt.Sprintf("MAZE_FILE_%d", req.MazeID))

	board, err := maze.ParseMaze(filename)
	if err != nil {
		utils.LogError(ctx, err, "failed to parse maze")
		http.Error(w, utils.Internal, http.StatusInternalServerError)
		return
	}

	if err = req.Validate(app.cfg, len(board[0]), len(board)); err != nil {
		utils.LogError(ctx, err, "failed to validate maze")
		http.Error(w, utils.Invalid, http.StatusBadRequest)
		return
	}

	newBoard, err := maze.UpdateMaze(filename, board, req.Points)
	if err != nil {
		utils.LogError(ctx, err, "failed to update maze")
		http.Error(w, utils.Internal, http.StatusInternalServerError)
		return
	}

	resp := models.UpdateMazeOutput{Map: toIntMap(newBoard)}
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		utils.LogError(ctx, err, utils.MsgErrMarshalResponse)
		http.Error(w, utils.Internal, http.StatusInternalServerError)
		return
	}
}

func (app *App) GetMazeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	mazeIDString := r.URL.Query().Get("labirint_id")
	mazeID, err := strconv.ParseInt(mazeIDString, 10, 64)
	if err != nil {
		utils.LogError(ctx, err, "failed to parse maze id")
		http.Error(w, utils.Invalid, http.StatusBadRequest)
		return
	}

	req := models.GetMazeInput{MazeID: int(mazeID)}
	if err = req.Validate(app.cfg); err != nil {
		utils.LogError(ctx, err, "invalid maze id")
		http.Error(w, utils.Invalid, http.StatusBadRequest)
		return
	}

	filename := os.Getenv(fmt.Sprintf("MAZE_FILE_%d", req.MazeID))

	board, err := maze.ParseMaze(filename)
	if err != nil {
		utils.LogError(ctx, err, "failed to parse maze")
		http.Error(w, utils.Internal, http.StatusInternalServerError)
		return
	}

	resp := models.GetMazeOutput{Map: toIntMap(board)}
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		utils.LogError(ctx, err, utils.MsgErrMarshalResponse)
		http.Error(w, utils.Internal, http.StatusInternalServerError)
		return
	}
}

func toIntMap(board [][]bool) [][]int {
	result := make([][]int, len(board))
	for i, row := range board {
		result[i] = make([]int, len(row))
		for j, cell := range row {
			if cell {
				result[i][j] = 1
			} else {
				result[i][j] = 0
			}
		}
	}

	return result
}
