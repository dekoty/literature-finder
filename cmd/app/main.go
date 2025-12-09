package main

import (
	"literature-finder/internal/adapter/persistence"
	"literature-finder/internal/adapter/server"
	"literature-finder/internal/config"
	"literature-finder/internal/usecase/search"
	"literature-finder/pkg/database"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Ошибка инициализации конфигурации: ", err)
	}

	dbKey, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Ошибка подключения к БД: ", err)
	}

	defer dbKey.Close()

	dbRepo := persistence.NewPostgresRepository(dbKey)

	repoGoogleBooks := persistence.NewGoogleBooksRepository(cfg.APIKey)
	repoOpenLibrary := persistence.NewOpenLibraryRepository()

	repo := persistence.NewMultiRepository(repoGoogleBooks, repoOpenLibrary)
	usecase := search.New(repo)
	handler := server.NewHandler(dbRepo, usecase)

	router := mux.NewRouter()

	router.HandleFunc("/", handler.Home).Methods("GET")
	router.HandleFunc("/search", handler.Search).Methods("GET")
	router.HandleFunc("/favorites", handler.FavoritesPageHandler).Methods("GET")
	router.HandleFunc("/save-favorite", handler.SaveFavoriteHandler).Methods("POST")
	router.HandleFunc("/delete-favorite", handler.DeleteFavoriteHandler).Methods("POST")
	router.HandleFunc("/clear-favorites", handler.ClearFavoritesHandler).Methods("POST")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	log.Printf("Сервер запущен на порту %s (http://localhost%s)", cfg.ServerAddress, cfg.ServerAddress)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
