package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"algo/algorithms"
	"algo/algorithms/a_star"
	"algo/handlers/models"
	"algo/maze"
	"algo/utils"
)

func SolveMazeHandler(w http.ResponseWriter, r *http.Request) {
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

	boundaryCells := make([][2]int, 0)
	if len(req.End) == 0 {
		boundaryCells = a_star.GetBoundaryCells(board, req.Start.X, req.Start.Y)
	} else {
		for _, end := range req.End {
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
