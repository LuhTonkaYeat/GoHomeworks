package http

import (
	"encoding/json"
	"net/http"
	"strings"

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
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	CreatedAt   string `json:"created_at"`
}

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type PingResponse struct {
	Status   string          `json:"status"`
	Services []ServiceStatus `json:"services"`
}

// GetRepository godoc
// @Summary Get repository information
// @Description Get information about a GitHub repository by URL
// @Tags repositories
// @Accept json
// @Produce json
// @Param url query string true "GitHub repository URL (e.g., https://github.com/golang/go)"
// @Success 200 {object} RepositoryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/repositories/info [get]
func (h *Handler) GetRepository(w http.ResponseWriter, r *http.Request) {
	urlParam := r.URL.Query().Get("url")
	if urlParam == "" {
		sendError(w, http.StatusBadRequest, "url parameter is required")
		return
	}

	owner, repo := parseGitHubURL(urlParam)
	if owner == "" || repo == "" {
		sendError(w, http.StatusBadRequest, "invalid github url")
		return
	}

	repository, err := h.repoUseCase.GetRepository(r.Context(), owner, repo)
	if err != nil {
		errMsg := err.Error()

		if strings.Contains(errMsg, "not found") {
			sendError(w, http.StatusNotFound, errMsg)
			return
		}

		sendError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response := RepositoryResponse{
		FullName:    repository.Name,
		Description: repository.Description,
		Stars:       repository.Stars,
		Forks:       repository.Forks,
		CreatedAt:   repository.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	sendJSON(w, http.StatusOK, response)
}

// Ping godoc
// @Summary Check services status
// @Description Ping processor and subscriber services
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} PingResponse
// @Failure 503 {object} PingResponse
// @Router /api/ping [get]
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	status, services, err := h.repoUseCase.Ping(r.Context())
	if err != nil {
		sendError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	svcStatuses := make([]ServiceStatus, len(services))
	for i, s := range services {
		svcStatuses[i] = ServiceStatus{
			Name:   s.Name,
			Status: s.Status,
		}
	}

	response := PingResponse{
		Status:   status,
		Services: svcStatuses,
	}

	if status != "ok" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	sendJSON(w, http.StatusOK, response)
}

func parseGitHubURL(url string) (string, string) {
	url = strings.TrimPrefix(url, "https://github.com/")
	url = strings.TrimPrefix(url, "http://github.com/")
	url = strings.TrimPrefix(url, "github.com/")
	url = strings.TrimSuffix(url, "/")

	parts := strings.Split(url, "/")
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}
	return "", ""
}

func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, ErrorResponse{Error: message})
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
