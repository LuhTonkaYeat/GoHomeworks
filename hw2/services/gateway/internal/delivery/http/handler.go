package http

import (
	"encoding/json"
	"net/http"

	"github.com/LuhTonkaYeat/GoHomeworks/hw2/services/gateway/internal/usecase"
)

type Handler struct {
	repoUseCase usecase.RepositoryUseCase
}

func NewHandler(repoUseCase usecase.RepositoryUseCase) *Handler {
	return &Handler{
		repoUseCase: repoUseCase,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type RepositoryResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	CreatedAt   string `json:"created_at"`
}

// GetRepository godoc
// @Summary Get repository information
// @Description Get information about a GitHub repository by owner and repo name
// @Tags repositories
// @Accept json
// @Produce json
// @Param owner query string true "Repository owner (username or organization)"
// @Param repo query string true "Repository name"
// @Success 200 {object} RepositoryResponse "Successfully retrieved repository information"
// @Failure 400 {object} ErrorResponse "Bad request - missing parameters"
// @Failure 404 {object} ErrorResponse "Repository not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /repo [get]
func (h *Handler) GetRepository(w http.ResponseWriter, r *http.Request) {
	owner := r.URL.Query().Get("owner")
	repo := r.URL.Query().Get("repo")

	if owner == "" || repo == "" {
		sendError(w, http.StatusBadRequest, "owner and repo are required")
		return
	}

	repository, err := h.repoUseCase.GetRepository(r.Context(), owner, repo)
	if err != nil {
		errMsg := err.Error()

		if len(errMsg) > 0 && (errMsg[:10] == "repository" ||
			(len(errMsg) > 10 && errMsg[len(errMsg)-10:] == "not found")) {
			sendError(w, http.StatusNotFound, errMsg)
			return
		}

		sendError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response := RepositoryResponse{
		Name:        repository.Name,
		Description: repository.Description,
		Stars:       repository.Stars,
		Forks:       repository.Forks,
		CreatedAt:   repository.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	sendJSON(w, http.StatusOK, response)
}

func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, ErrorResponse{Error: message})
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
