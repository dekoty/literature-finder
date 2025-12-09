package persistence

import (
	"literature-finder/internal/module/literature"
	"log"
	"sync"
)

type MultiRepository struct {
	repos []literature.Repository
}

func NewMultiRepository(repos ...literature.Repository) *MultiRepository {
	return &MultiRepository{repos: repos}
}

func (m *MultiRepository) Search(quary string) ([]literature.Literature, error) {
	var all []literature.Literature
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, r := range m.repos {
		wg.Add(1)
		go func(r literature.Repository) {
			defer wg.Done()
			results, err := r.Search(quary)

			if err != nil {
				log.Printf("Внимание: Ошибка при поиске в репозитории %T: %v", r, err)
				return
			}

			mu.Lock()
			all = append(all, results...)
			mu.Unlock()
		}(r)
	}

	wg.Wait()

	return all, nil
}
