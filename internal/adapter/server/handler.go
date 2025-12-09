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

	tmpl := template.New("base").Funcs(funcMap)

	globPattern := filepath.Join("templates", "*.html")

	tmpl = template.Must(tmpl.ParseGlob(globPattern))

	return &Handler{
		SearchUC:  uc,
		DBRepo:    dbRepo,
		Templates: tmpl,
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {

	if err := h.Templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Printf("Ошибка рендеринга шаблона index.html: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	query = strings.TrimSpace(query)

	if query == "" {
		h.Home(w, r)
		return
	}

	result, err := h.SearchUC.SearchLiterature(query)
	if err != nil {
		log.Printf("Ошибка поиска для запроса '%s': %v", query, err)
		http.Error(w, "Не удалось выполнить поиск", http.StatusInternalServerError)
		return
	}

	data := struct {
		Results []literature.Literature
	}{
		Results: result,
	}

	if err := h.Templates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("Ошибка рендеринга результатов поиска: %v", err)
		http.Error(w, "Ошибка отображения данных", http.StatusInternalServerError)
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
		log.Printf("Ошибка парсинга формы сохранения: %v", err)
		http.Error(w, "Некорректные данные формы", http.StatusBadRequest)
		return
	}

	authorsString := r.FormValue("authors")
	book := literature.Literature{
		ID:        strings.TrimSpace(r.FormValue("book_id")),
		Title:     strings.TrimSpace(r.FormValue("title")),
		Authors:   strings.Split(authorsString, "; "),
		Thumbnail: r.FormValue("thumbnail"),
		Link:      r.FormValue("link"),
		Status:    "favorite",
		Year:      r.FormValue("year"),
	}

	if book.ID == "" || book.Title == "" {
		log.Println("Попытка сохранить книгу без ID или названия")
		http.Error(w, "Отсутствуют обязательные данные книги", http.StatusBadRequest)
		return
	}

	if err := h.DBRepo.SaveBook(userID, book); err != nil {
		log.Printf("Ошибка сохранения книги (ID: %s) в БД: %v", book.ID, err)
		http.Error(w, "Внутренняя ошибка базы данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Книга успешно добавлена в избранное",
		"status":  "success",
	}); err != nil {
		log.Printf("Ошибка записи JSON ответа: %v", err)
	}
}

func (h *Handler) FavoritesPageHandler(w http.ResponseWriter, r *http.Request) {
	userID := util.GetUserID(w, r)

	books, err := h.DBRepo.GetBooksByUserID(userID, "favorite")
	if err != nil {
		log.Printf("Ошибка получения избранного для user %s: %v", userID, err)
		http.Error(w, "Не удалось загрузить библиотеку", http.StatusInternalServerError)
		return
	}

	data := struct {
		Results []literature.Literature
	}{
		Results: books,
	}

	if err := h.Templates.ExecuteTemplate(w, "favorites.html", data); err != nil {
		log.Printf("Ошибка рендеринга favorites.html: %v", err)
		http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
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
		log.Printf("Ошибка парсинга формы удаления: %v", err)
		http.Error(w, "Ошибка обработки запроса", http.StatusBadRequest)
		return
	}

	bookID := strings.TrimSpace(r.FormValue("book_id"))
	if bookID == "" {
		http.Error(w, "Отсутствует ID книги", http.StatusBadRequest)
		return
	}

	if err := h.DBRepo.DeleteBook(userID, bookID); err != nil {
		log.Printf("Ошибка удаления книги %s из БД: %v", bookID, err)
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
		log.Printf("Ошибка полной очистки избранного для user %s: %v", userID, err)
		http.Error(w, "Не удалось очистить список", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/favorites", http.StatusSeeOther)
}
