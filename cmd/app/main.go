package main

import (
	"database/sql"
	"fmt"
	"literature-finder/internal/adapter/persistence"
	"literature-finder/internal/adapter/server"
	"literature-finder/internal/usecase/search"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Print(".env не загрузился\n")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Print("GOOGLE_API_KEY is not set\n")
	}

	dbKey, err := sql.Open("postgres", os.Getenv("DATASOURCENAME"))

	if err != nil {
		log.Fatal("Ошибка открытия соединения с базой данных: ", err)
	}

	if err := dbKey.Ping(); err != nil {

		log.Fatalf("Ошибка при пинге базы данных: %v", err)
	}

	dbRepo := persistence.NewPostgresRepository(dbKey)

	repoGoogleBooks := persistence.NewGoogleBooksRepository(apiKey)
	repoOpenLibrary := persistence.NewOpenLibraryRepository()

	repo := persistence.NewMultiRepository(repoGoogleBooks, repoOpenLibrary)

	usecase := search.New(repo)
	handler := server.NewHandler(dbRepo, usecase)

	router := mux.NewRouter()

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	router.HandleFunc("/", handler.Home)
	router.HandleFunc("/search", handler.Search)
	router.HandleFunc("/favorites", handler.FavoritesPageHandler)
	router.HandleFunc("/save-favorite", handler.SaveFavoriteHandler)
	router.HandleFunc("/delete-favorite", handler.DeleteFavoriteHandler)
	router.HandleFunc("/clear-favorites", handler.ClearFavoritesHandler)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	fmt.Println("Сервер запущен на порту 8080...")

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Error %s\n", err)

	}

}
