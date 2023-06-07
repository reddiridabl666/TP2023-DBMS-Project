package usecase

import (
	"forum/internal/pkg/domain"
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

func (u *PostUsecase) AddPosts(thread *domain.Thread, posts domain.PostBatch) error {
	return u.posts.AddPosts(thread, posts)
}

func (u *PostUsecase) GetPost(id int64) (*domain.Post, error) {
	return u.posts.GetPost(id)
}

func (u *PostUsecase) Update(post *domain.Post) error {
	return u.posts.Update(post)
}

func (u *PostUsecase) GetPosts(params *domain.PostListParams) (domain.PostBatch, error) {
	switch params.Sort {
	case domain.SortFlat:
		return u.posts.GetPostsFlat(params)
	case domain.SortTree:
		return u.posts.GetPostsTree(params)
	case domain.SortParent:
		return u.posts.GetPostsParent(params)
	}
	return nil, domain.ErrInvalidArgument
}
