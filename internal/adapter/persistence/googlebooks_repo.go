package persistence

import (
	"encoding/json"
	"fmt"
	"literature-finder/internal/module/literature"
	"net/http"
	"net/url"
)

type GoogleBooksRepository struct {
	APIKey string
}

func NewGoogleBooksRepository(apiKey string) *GoogleBooksRepository {
	return &GoogleBooksRepository{APIKey: apiKey}
}

func (a *GoogleBooksRepository) Search(quary string) ([]literature.Literature, error) {
	urll := "https://www.googleapis.com/books/v1/volumes?q=%s&fields=items(volumeInfo(title,authors,publishedDate,description,infoLink,imageLinks/thumbnail))&key=%s"
	resp, err := http.Get(fmt.Sprintf(urll, url.QueryEscape(quary), a.APIKey))

	if err != nil {
		return nil, err
	}

	var apiResp googleBooksResponse

	err = json.NewDecoder(resp.Body).Decode(&apiResp)

	if err != nil {
		return nil, err
	}

	var results []literature.Literature

	for _, b := range apiResp.Items {
		v := b.VolumeInfo

		book := literature.Literature{
			Title:     v.Title,
			Authors:   v.Authors,
			Year:      v.PublishedDate,
			Thumbnail: v.ImageLinks.Thumbnail,
			Link:      v.InfoLink,
		}

		results = append(results, book)

	}

	return results, nil
}
