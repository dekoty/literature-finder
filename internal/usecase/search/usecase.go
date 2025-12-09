package search

import (
	"fmt"
	"literature-finder/internal/module/literature"
)

type UseCase struct {
	repo literature.Repository
}

func New(repo literature.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) SearchLiterature(query string) ([]literature.Literature, error) {
	results, err := uc.repo.Search(query)
	if err != nil {

		return nil, fmt.Errorf("ошибка бизнес-логики при поиске по запросу '%s': %w", query, err)
	}
	return results, nil
}
