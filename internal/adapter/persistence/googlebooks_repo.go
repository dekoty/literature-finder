package persistence

import (
	"encoding/json"
	"fmt"
	"literature-finder/internal/module/literature"
	"log"
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
	urll := "https://www.googleapis.com/books/v1/volumes?q=%s&maxResults=40&projection=full&printType=books&fields=items(id,volumeInfo/title,volumeInfo/authors,volumeInfo/publishedDate,volumeInfo/description,volumeInfo/infoLink,volumeInfo/imageLinks/thumbnail)&key=%s"
	resp, err := http.Get(fmt.Sprintf(urll, url.QueryEscape(quary), a.APIKey))

	if err != nil {
		return nil, fmt.Errorf("ошибка сетевого запроса к Google Books: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google Books API вернул статус %d. Проверьте API ключ или лимиты", resp.StatusCode)
	}

	var apiResp googleBooksResponse

	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования JSON от Google Books: %w", err)
	}

	var results []literature.Literature

	for _, b := range apiResp.Items {
		v := b.VolumeInfo
		var thumbnailURL string
		if v.ImageLinks.Thumbnail != "" {
			thumbnailURL = v.ImageLinks.Thumbnail
		} else {
			thumbnailURL = "/static/images/upscaled-image.png"
		}

		book := literature.Literature{
			ID:        b.ID,
			Title:     v.Title,
			Authors:   v.Authors,
			Year:      v.PublishedDate,
			Thumbnail: thumbnailURL,
			Link:      v.InfoLink,
		}

		log.Printf("Mapper Output: ID=%s, Title=%s", book.ID, book.Title)

		results = append(results, book)
	}

	return results, nil
}
