package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"literature-finder/internal/adapter/persistence"
	"literature-finder/internal/adapter/server"
	"literature-finder/internal/usecase/search"
	"net/http"
	"os"
)

func main() {

	err := godotenv.Load("../../.env")

	if err != nil {
		fmt.Print(".env не загрузился\n")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Print("GOOGLE_API_KEY is not set\n")
	}

	repoGoogleBooks := persistence.NewGoogleBooksRepository(apiKey)
	repoOpenLibrary := persistence.NewOpenLibraryRepository()

	repo := persistence.NewMultiRepository(repoGoogleBooks, repoOpenLibrary)

	usecase := search.New(repo)
	handler := server.NewHandler(usecase)

	router := mux.NewRouter()

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	router.HandleFunc("/search", handler.Search)

	fmt.Println("Сервер запущен на порту 8080...")

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Error %s\n", err)

	}

}
