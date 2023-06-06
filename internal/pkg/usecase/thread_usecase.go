package usecase

import (
	"strconv"

	"forum/internal/pkg/domain"
	"forum/internal/pkg/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ThreadUsecase struct {
	forums  *repository.ForumRepository
	threads *repository.ThreadRepository
	users   *repository.UserRepository
}

func NewThreadUsecase(db *pgxpool.Pool) *ThreadUsecase {
	return &ThreadUsecase{
		forums:  repository.NewForumRepository(db),
		threads: repository.NewThreadRepository(db),
		users:   repository.NewUserRepository(db),
	}
}

func (u *ThreadUsecase) Create(thread *domain.Thread) error {
	user, err := u.users.GetByNickname(thread.Author)
	if err != nil {
		return err
	}

	forum, err := u.forums.Get(thread.Forum)
	if err != nil {
		return err
	}

	thread.Author = user.Nickname
	thread.Forum = forum.Slug
	return u.threads.Create(forum.Id, user.Id, thread)
}

func (u *ThreadUsecase) Get(slugOrId string) (*domain.Thread, error) {
	id, err := strconv.Atoi(slugOrId)
	if err == nil {
		return u.threads.GetById(id)
	}
	return u.threads.GetBySlug(slugOrId)
}

func (u *ThreadUsecase) GetById(id int) (*domain.Thread, error) {
	return u.threads.GetById(id)
}

func (u *ThreadUsecase) Update(slugOrId string, thread *domain.Thread) error {
	previous, err := u.Get(slugOrId)
	if err != nil {
		return err
	}

	if thread.Title == "" {
		thread.Title = previous.Title
	}

	if thread.Message == "" {
		thread.Message = previous.Message
	}

	thread.Author = previous.Author
	thread.Forum = previous.Forum
	thread.Created = previous.Created
	thread.Id = previous.Id
	thread.Slug = previous.Slug

	return u.threads.Update(thread)
}

func (u *ThreadUsecase) GetByForum(params *domain.ThreadListParams) (domain.ThreadBatch, error) {
	forum, err := u.forums.Get(params.Forum)
	if err != nil {
		return nil, err
	}
	params.ForumId = forum.Id

	if params.Limit < 1 {
		params.Limit = 100
	}
	return u.threads.GetByForum(params)
}
