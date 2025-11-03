package search

import "literature-finder/internal/module/literature"

type UseCase struct {
	repo literature.Repository
}

func New(repo literature.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) SearchLiterature(query string) ([]literature.Literature, error) {

	return uc.repo.Search(query)
}
