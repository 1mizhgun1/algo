package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"algo/algorithms"
	"algo/algorithms/a_star"
	"algo/algorithms/lazy_theta_star"
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
		boundaryCells = algorithms.GetBoundaryCells(board, req.Start.X, req.Start.Y)
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
		result       models.SolveMazeOutput
		distance     int
		shortestPath []algorithms.Node
	)

	startTime := time.Now()
	switch req.AlgorithmID {
	case 1:
		distance, shortestPath = a_star.AStar(board, req.Start.X, req.Start.Y, boundaryCells)
	case 2:
		distance, shortestPath = lazy_theta_star.LazyThetaStar(board, req.Start.X, req.Start.Y, boundaryCells)
	default:
		utils.LogErrorMessage(ctx, fmt.Sprintf("invalid algorithm id=%d", req.AlgorithmID))
		http.Error(w, utils.Invalid, http.StatusBadRequest)
		return
	}
	endTime := time.Now()

	if distance == algorithms.PathNotFound {
		if err = json.NewEncoder(w).Encode(result); err != nil {
			utils.LogError(ctx, err, utils.MsgErrMarshalResponse)
			http.Error(w, utils.Internal, http.StatusInternalServerError)
			return
		}
		return
	}

	if len(shortestPath) == 0 {
		utils.LogErrorMessage(ctx, "path is empty")
		http.Error(w, utils.Internal, http.StatusInternalServerError)
		return
	}

	result.Path = make([]models.Tranzition, len(shortestPath)-1)
	result.Dist = distance
	result.ExecutionTime = endTime.Sub(startTime)

	for i := 1; i < len(shortestPath); i++ {
		result.Path[i-1] = models.Tranzition{
			Start: models.Point{X: shortestPath[i-1].X, Y: shortestPath[i-1].Y},
			End:   models.Point{X: shortestPath[i].X, Y: shortestPath[i].Y},
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

func (app *App) RestoreMazeHandler(w http.ResponseWriter, r *http.Request) {
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
	originalFilename := getOriginalFilename(filename)

	board, err := maze.RestoreMaze(filename, originalFilename)
	if err != nil {
		utils.LogError(ctx, err, "failed to update maze")
		http.Error(w, utils.Internal, http.StatusInternalServerError)
		return
	}

	resp := models.UpdateMazeOutput{Map: toIntMap(board)}
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

func getOriginalFilename(filename string) string {
	ext := path.Ext(filename)
	return strings.TrimSuffix(filename, ext) + "_default" + ext
}
