package persistence

type googleBooksResponse struct {
	Items []struct {
		VolumeInfo struct {
			Title         string   `json:"title"`
			Authors       []string `json:"authors"`
			PublishedDate string   `json:"publishedDate"`
			Description   string   `json:"description"`
			ImageLinks    struct {
				Thumbnail string `json:"thumbnail"`
			} `json:"imageLinks"`
			InfoLink string `json:"infoLink"`
		} `json:"volumeInfo"`
	} `json:"items"`
}
