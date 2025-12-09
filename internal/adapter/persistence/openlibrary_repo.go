package persistence

import (
	"encoding/json"
	"fmt"
	"literature-finder/internal/module/literature"
	"net/http"
	"net/url"
)

type OpenLibraryRepository struct{}

type openLibraryResponse struct {
	Docs []struct {
		Title            string   `json:"title"`
		AuthorName       []string `json:"author_name"`
		FirstPublishYear int      `json:"first_publish_year"`
		CoverI           int      `json:"cover_i"`
		Key              string   `json:"key"`
	} `json:"docs"`
}

func NewOpenLibraryRepository() *OpenLibraryRepository {
	return &OpenLibraryRepository{}
}

func (a *OpenLibraryRepository) Search(quary string) ([]literature.Literature, error) {
	urll := "https://openlibrary.org/search.json?q=%s&lang=ru"
	resp, err := http.Get(fmt.Sprintf(urll, url.QueryEscape(quary)))

	if err != nil {
		return nil, err
	}

	var apiResp openLibraryResponse

	err = json.NewDecoder(resp.Body).Decode(&apiResp)

	if err != nil {
		return nil, err
	}

	var results []literature.Literature

	for _, b := range apiResp.Docs {
		thumb := ""
		if b.CoverI != 0 {
			thumb = fmt.Sprintf("https://covers.openlibrary.org/b/id/%d-M.jpg", b.CoverI)
		}

		if thumb == "" {
			thumb = "/static/images/upscaled-image.png"

		}

		book := literature.Literature{
			ID:        b.Key,
			Title:     b.Title,
			Authors:   b.AuthorName,
			Year:      fmt.Sprintf("%d", b.FirstPublishYear),
			Thumbnail: thumb,
			Link:      fmt.Sprintf("https://openlibrary.org%s", b.Key),
		}

		results = append(results, book)

	}

	return results, nil

}
