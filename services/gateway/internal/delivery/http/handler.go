package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/LuhTonkaYeat/GoHomeworks/services/gateway/internal/usecase"
)

type Handler struct {
	repoUseCase         usecase.RepositoryUseCase
	subscriptionUseCase usecase.SubscriptionUseCase
}

func NewHandler(repoUseCase usecase.RepositoryUseCase, subscriptionUseCase usecase.SubscriptionUseCase) *Handler {
	return &Handler{
		repoUseCase:         repoUseCase,
		subscriptionUseCase: subscriptionUseCase,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type RepositoryResponse struct {
	FullName    string
	Description string
	Stars       int
	Forks       int
	CreatedAt   string
}

type ServiceStatus struct {
	Name   string
	Status string
}

type PingResponse struct {
	Status   string
	Services []ServiceStatus
}

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

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Owner string `json:"owner"`
		Repo  string `json:"repo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Owner == "" || req.Repo == "" {
		sendError(w, http.StatusBadRequest, "owner and repo are required")
		return
	}

	userID := "default"

	err := h.subscriptionUseCase.CreateSubscription(r.Context(), req.Owner, req.Repo, userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusCreated, map[string]string{"message": "subscription created"})
}

func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")

	if owner == "" || repo == "" {
		sendError(w, http.StatusBadRequest, "owner and repo are required")
		return
	}

	userID := "default"

	err := h.subscriptionUseCase.DeleteSubscription(r.Context(), owner, repo, userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"message": "subscription deleted"})
}

func (h *Handler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID := "default"

	subscriptions, err := h.subscriptionUseCase.GetSubscriptions(r.Context(), userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, subscriptions)
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

func (h *Handler) GetSubscriptionsInfo(w http.ResponseWriter, r *http.Request) {
	userID := "default"

	repositories, err := h.subscriptionUseCase.GetSubscriptionsInfo(r.Context(), userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"repositories": repositories,
	})
}
