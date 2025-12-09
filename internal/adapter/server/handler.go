package server

import (
	"encoding/json"
	"html/template"
	"literature-finder/internal/adapter/persistence"
	"literature-finder/internal/module/literature"
	"literature-finder/internal/usecase/search"
	"literature-finder/internal/util"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type Handler struct {
	SearchUC  *search.UseCase
	DBRepo    *persistence.PostgresRepository
	Templates *template.Template
}

func NewHandler(dbRepo *persistence.PostgresRepository, uc *search.UseCase) *Handler {

	funcMap := template.FuncMap{
		"join": strings.Join,
	}

	tmpl := template.New("index.html").Funcs(funcMap)

	globPattern := filepath.Join("templates", "*.html")

	tmpl = template.Must(tmpl.ParseGlob(globPattern))

	return &Handler{
		SearchUC:  uc,
		DBRepo:    dbRepo,
		Templates: tmpl,
	}

}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {

	if err := h.Templates.Execute(w, nil); err != nil {
		http.Error(w, "Ошибка загрузки шаблона: "+err.Error(), 500)
		return
	}

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

	data := struct {
		Results []literature.Literature
	}{
		Results: result,
	}

	err = h.Templates.Execute(w, data)

	if err != nil {
		return
	}

}

func (h *Handler) SaveFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	userID := util.GetUserID(w, r)

	if err := r.ParseForm(); err != nil {
		log.Printf("Ошибка парсинга формы: %v", err)
		http.Error(w, "Ошибка парсинга запроса", http.StatusBadRequest)
		return
	}

	log.Printf("Form Received: book_id='%s', title='%s'",
		r.FormValue("book_id"),
		r.FormValue("title"),
	)

	authorsString := r.FormValue("authors")

	book := literature.Literature{
		ID:    strings.TrimSpace(r.FormValue("book_id")),
		Title: strings.TrimSpace(r.FormValue("title")),

		Authors:   strings.Split(authorsString, "; "),
		Thumbnail: r.FormValue("thumbnail"),
		Link:      r.FormValue("link"),
		Status:    "favorite",
		Year:      r.FormValue("year"),
	}

	if book.ID == "" || book.Title == "" {
		http.Error(w, "Отсутствуют обязательные данные книги", http.StatusBadRequest)
		return
	}

	if err := h.DBRepo.SaveBook(userID, book); err != nil {
		log.Printf("Ошибка сохранения книги в БД для user %s: %v", userID, err)
		http.Error(w, "Не удалось сохранить книгу", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Книга успешно добавлена в избранное",
		"status":  "success",
	})
}

func (h *Handler) FavoritesPageHandler(w http.ResponseWriter, r *http.Request) {

	userID := util.GetUserID(w, r)

	books, err := h.DBRepo.GetBooksByUserID(userID, "favorite")
	if err != nil {
		log.Printf("Ошибка получения избранных книг для user %s: %v", userID, err)
		http.Error(w, "Не удалось загрузить список книг.", http.StatusInternalServerError)
		return
	}

	data := struct {
		Results []literature.Literature
	}{
		Results: books,
	}

	err = h.Templates.ExecuteTemplate(w, "favorites.html", data)
	if err != nil {
		log.Printf("Ошибка выполнения шаблона favorites.html: %v", err)
		return
	}
}

func (h *Handler) DeleteFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	userID := util.GetUserID(w, r)

	if err := r.ParseForm(); err != nil {
		log.Printf("Ошибка парсинга формы: %v", err)
		http.Error(w, "Ошибка парсинга запроса", http.StatusBadRequest)
		return
	}

	bookID := strings.TrimSpace(r.FormValue("book_id"))

	if bookID == "" {
		http.Error(w, "Отсутствует ID книги", http.StatusBadRequest)
		return
	}

	if err := h.DBRepo.DeleteBook(userID, bookID); err != nil {
		log.Printf("Ошибка удаления книги из БД: %v", err)
		http.Error(w, "Не удалось удалить книгу", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/favorites", http.StatusSeeOther)

}

func (h *Handler) ClearFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	userID := util.GetUserID(w, r)

	if err := h.DBRepo.ClearFavorites(userID); err != nil {
		log.Printf("Ошибка очистки избранного в БД: %v", err)
		http.Error(w, "Не удалось очистить список", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/favorites", http.StatusSeeOther)
}
