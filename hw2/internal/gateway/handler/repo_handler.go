package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/LuhTonkaYeat/GoHomeworks/hw2/internal/dto"
	"github.com/LuhTonkaYeat/GoHomeworks/hw2/internal/gateway/client"
)

type RepoHandler struct {
	collectorClient *client.CollectorClient
}

func NewRepoHandler(collectorClient *client.CollectorClient) *RepoHandler {
	return &RepoHandler{
		collectorClient: collectorClient,
	}
}

func (h *RepoHandler) GetRepository(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	owner := r.URL.Query().Get("owner")
	repo := r.URL.Query().Get("repo")

	if owner == "" || repo == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Error:   "Bad Request",
			Code:    http.StatusBadRequest,
			Message: "Both 'owner' and 'repo' parameters are required",
		})
		return
	}

	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)
	repo = strings.TrimPrefix(repo, "/")
	repo = strings.TrimSuffix(repo, "/")

	resp, err := h.collectorClient.GetRepository(r.Context(), owner, repo)
	if err != nil {
		statusCode, message := client.MapGrpcErrorToHTTP(err)
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Error:   http.StatusText(statusCode),
			Code:    statusCode,
			Message: message,
		})
		return
	}

	response := dto.RepoResponse{
		Name:        resp.Name,
		Description: resp.Description,
		Stars:       int(resp.Stars),
		Forks:       int(resp.Forks),
		CreatedAt:   resp.CreatedAt,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
