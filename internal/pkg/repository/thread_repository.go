package repository

import (
	"context"
	"errors"

	"forum/internal/pkg/domain"
	"forum/internal/pkg/utils"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ThreadRepository struct {
	db *pgxpool.Pool
}

func NewThreadRepository(db *pgxpool.Pool) *ThreadRepository {
	return &ThreadRepository{
		db: db,
	}
}

func (repo *ThreadRepository) Create(forumId, authorId int, thread *domain.Thread) error {
	err := repo.db.QueryRow(context.Background(),
		`INSERT INTO thread(forum_id, author_id, title, message, slug, created_at)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		forumId, authorId, thread.Title, thread.Message,
		thread.Slug, thread.Created.UTC().UnixNano(),
	).
		Scan(&thread.Id)

	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != pgerrcode.UniqueViolation {
		return err
	}

	err = repo.db.QueryRow(context.Background(),
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

	err := repo.db.QueryRow(context.Background(),
		`SELECT t.id, t.title, u.nickname, f.slug, t.forum_id,
			 	t.message, t.rating, t.slug, t.created_at
		 FROM thread t JOIN users u ON t.author_id = u.id
		 			   JOIN forum f ON t.forum_id  = f.id
		 WHERE t.id = $1`, id).
		Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.ForumId,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created,
		)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (repo *ThreadRepository) GetBySlug(slug string) (*domain.Thread, error) {
	thread := &domain.Thread{}

	err := repo.db.QueryRow(context.Background(),
		`SELECT t.id, t.title, u.nickname, f.slug, t.forum_id,
			 	t.message, t.rating, t.slug, t.created_at
		 FROM thread t JOIN users u ON t.author_id = u.id
		 			   JOIN forum f ON t.forum_id  = f.id
		 WHERE lower(t.slug) = lower($1)`, slug).
		Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.ForumId,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created,
		)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (repo *ThreadRepository) Update(thread *domain.Thread) error {
	_, err := repo.db.Exec(context.Background(),
		`UPDATE thread SET title = $1, message = $2 WHERE id = $3`,
		thread.Title, thread.Message, thread.Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (repo *ThreadRepository) GetByForum(params *domain.ThreadListParams) (domain.ThreadBatch, error) {
	query := `SELECT t.id, t.title, u.nickname, f.slug,
					 t.message, t.rating, t.slug, t.created_at
				FROM thread t JOIN users u ON t.author_id = u.id
							  JOIN forum f ON t.forum_id  = f.id
				WHERE t.forum_id = $1 AND t.created_at`

	if !params.Desc {
		query += " >= $2 ORDER BY t.created_at"
	} else {
		if params.Since == 0 {
			params.Since = utils.MaxTime.UnixNano()
		}
		query += " <= $2 ORDER BY t.created_at DESC"
	}
	query += " LIMIT $3"

	rows, err := repo.db.Query(context.Background(), query, params.ForumId, params.Since, params.Limit)
	res := []*domain.Thread{}

	if err == pgx.ErrNoRows {
		return res, nil
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		thread := &domain.Thread{}
		err = rows.Scan(
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
			return nil, err
		}
		res = append(res, thread)
	}
	return res, nil
}
