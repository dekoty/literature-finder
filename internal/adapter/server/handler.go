package server

import (
	"html/template"
	"literature-finder/internal/module/literature"
	"literature-finder/internal/usecase/search"
	"net/http"
	"path/filepath"
	"strings"
)

type Handler struct {
	SearchUC *search.UseCase
}

func NewHandler(uc *search.UseCase) *Handler {
	return &Handler{SearchUC: uc}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона: "+err.Error(), 500)
		return
	}

	tmpl.Execute(w, nil)
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	quary := r.URL.Query().Get("q")

	quary = strings.TrimSpace(quary)
	if quary == "" {
		h.Home(w, r)
		return
	}

	result, err := h.SearchUC.SearchLiterature(quary)

	if err != nil {
		return
	}

	main := filepath.Join("templates", "index.html")

	tmpl, err := template.ParseFiles(main)

	if err != nil {
		return
	}

	data := struct {
		Results []literature.Literature
	}{
		Results: result,
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		return
	}

}
