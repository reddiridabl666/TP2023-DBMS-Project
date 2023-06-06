package usecase

import (
	"forum/internal/pkg/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostUsecase struct {
	posts   *repository.PostRepository
	threads *repository.ThreadRepository
}

func NewPostUsecase(db *pgxpool.Pool) *PostUsecase {
	return &PostUsecase{
		posts:   repository.NewPostRepository(db),
		threads: repository.NewThreadRepository(db),
	}
}
