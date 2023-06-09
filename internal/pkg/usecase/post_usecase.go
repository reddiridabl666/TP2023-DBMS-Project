package usecase

import (
	"forum/internal/pkg/domain"
	"forum/internal/pkg/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostUsecase struct {
	posts *repository.PostRepository
}

func NewPostUsecase(db *pgxpool.Pool) *PostUsecase {
	return &PostUsecase{
		posts: repository.NewPostRepository(db),
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

func (u *PostUsecase) GetPosts(thread *domain.Thread, params *domain.PostListParams) (domain.PostBatch, error) {
	var posts domain.PostBatch
	var err error

	switch params.Sort {
	case domain.SortFlat:
		posts, err = u.posts.GetPostsFlat(params)
	case domain.SortTree:
		posts, err = u.posts.GetPostsTree(params)
	case domain.SortParent:
		posts, err = u.posts.GetPostsParent(params)
	default:
		return nil, domain.ErrInvalidArgument
	}

	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		post.Forum = thread.Forum
	}
	return posts, nil
}
