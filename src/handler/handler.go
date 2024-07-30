package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	useCase "youtube_tracker/src/usecase/video"

	"github.com/go-chi/render"
)

type Handler struct {
	videouseCase useCase.VideoService
}

func NewHandler(videouseCase useCase.VideoService) *Handler {
	return &Handler{videouseCase}
}

func (h *Handler) SearchVideos(w http.ResponseWriter, r *http.Request) {
	keywords := r.URL.Query().Get("keyword")
	if keywords == "" {
		http.Error(w, "keyword is required", http.StatusBadRequest)
		return
	}
	keywordList := strings.Split(keywords, ",")
	if err := h.videouseCase.SearchVideos(r.Context(), keywordList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetVideos(w http.ResponseWriter, r *http.Request) {
	keywords := r.URL.Query().Get("keyword")
	if keywords == "" {
		http.Error(w, "keyword is required", http.StatusBadRequest)
		return
	}
	keywordList := strings.Split(keywords, ",")
	videos, err := h.videouseCase.GetVideos(r.Context(), keywordList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, videos)
}

func (h *Handler) GetPopularityScores(w http.ResponseWriter, r *http.Request) {
	keywords := r.URL.Query().Get("keyword")
	if keywords == "" {
		http.Error(w, "keyword is required", http.StatusBadRequest)
		return
	}
	keywordList := strings.Split(keywords, ",")

	scores, err := h.videouseCase.CalculatePopularityScore(keywordList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scores)
}
