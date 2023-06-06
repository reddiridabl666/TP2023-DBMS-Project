package repository

import (
	"database/sql"

	"forum/internal/pkg/domain"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

type ThreadRepository struct {
	db *sql.DB
}

func NewThreadRepository(db *sql.DB) *ThreadRepository {
	return &ThreadRepository{
		db: db,
	}
}

func (repo *ThreadRepository) Create(forumId, authorId int, thread *domain.Thread) error {
	err := repo.db.QueryRow(
		`INSERT INTO thread(forum_id, author_id, title, message, slug, created_at)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		forumId, authorId, thread.Title, thread.Message, thread.Slug, thread.Created).
		Scan(&thread.Id)

	if err == nil {
		return nil
	}

	pgError, ok := err.(pgx.PgError)
	if !ok || pgError.Code != pgerrcode.UniqueViolation {
		return err
	}

	err = repo.db.QueryRow(
		`SELECT t.id, t.title, u.nickname, f.slug,
			 	t.message, t.rating, t.slug, t.created_at
		 FROM thread t JOIN users u ON t.author_id = u.id
		 			   JOIN forum f ON t.forum_id  = f.id
		 WHERE lower(t.slug) = lower($1)`, thread.Slug).
		Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created,
		)
	if err != nil {
		return err
	}
	return domain.ErrAlreadyExists
}

func (repo *ThreadRepository) GetById(id int) (*domain.Thread, error) {
	thread := &domain.Thread{}

	err := repo.db.QueryRow(
		`SELECT t.id, t.title, u.nickname, f.slug,
			 	t.message, t.rating, t.slug, t.created_at
		 FROM thread t JOIN users u ON t.author_id = u.id
		 			   JOIN forum f ON t.forum_id  = f.id
		 WHERE t.id = $1`, id).
		Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created,
		)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (repo *ThreadRepository) GetBySlug(slug string) (*domain.Thread, error) {
	thread := &domain.Thread{}

	err := repo.db.QueryRow(
		`SELECT t.id, t.title, u.nickname, f.slug,
			 	t.message, t.votes, t.slug, t.created_at
		 FROM thread t JOIN users u ON t.author_id = u.id
		 			   JOIN forum f ON t.forum_id  = f.id
		 WHERE lower(t.slug) = lower($1)`, slug).
		Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created,
		)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return thread, nil
}
