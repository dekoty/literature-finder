package literature

type Repository interface {
	Search(quary string) ([]Literature, error)
}
