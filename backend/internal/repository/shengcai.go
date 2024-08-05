package repository

type ShengCaiRepository interface {
}

func NewShengCaiRepository(
	r *Repository,
) ShengCaiRepository {
	return &shengCaiRepository{
		Repository: r,
	}
}

type shengCaiRepository struct {
	*Repository
}
