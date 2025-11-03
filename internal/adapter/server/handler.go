package server

import (
	"encoding/json"
	"literature-finder/internal/usecase/search"
	"net/http"
)

type Handler struct {
	SearchUC *search.UseCase
}

func NewHandler(uc *search.UseCase) *Handler {
	return &Handler{SearchUC: uc}
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	quary := r.URL.Query().Get("q")

	result, err := h.SearchUC.SearchLiterature(quary)

	if err != nil {
		return
	}

	json.NewEncoder(w).Encode(result)

}
