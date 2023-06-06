package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"forum/internal/pkg/domain"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository struct {
	db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (repo *PostRepository) AddPosts(thread *domain.Thread, posts domain.PostBatch) error {
	ids, err := repo.getAuthorIds(posts)
	if err != nil {
		return err
	}

	createdAt := time.Now().UTC()
	query := `INSERT INTO post(thread_id, author_id, parent_id, message, created_at) VALUES `

	args := make([]interface{}, 0, len(posts)*3+2)
	args = append(args, thread.Id, createdAt)

	i := 3
	for _, post := range posts {
		post.Thread = thread.Id
		post.Forum = thread.Forum
		post.Created = createdAt

		query += fmt.Sprintf("($1, $%d, $%d, $%d, $2),", i, i+1, i+2)
		args = append(args, ids[post.Author], post.Parent, post.Message)
		i += 3
	}

	query = query[:len(query)-1] + " RETURNING id"

	rows, err := repo.db.Query(context.Background(), query, args...)
	if err != nil {
		return err
	}

	postIds, err := pgx.CollectRows[int64](rows, pgx.RowTo[int64])
	if err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return err
		}

		switch pgErr.Code {
		case pgerrcode.ForeignKeyViolation:
			return domain.ErrNoParent
		case pgerrcode.IntegrityConstraintViolation:
			return domain.ErrInvalidParent
		}
	}

	for i, id := range postIds {
		posts[i].Id = id
	}

	return nil
}

func (repo *PostRepository) getAuthorIds(posts domain.PostBatch) (map[string]int, error) {
	res := make(map[string]int, len(posts))

	for _, post := range posts {
		_, ok := res[post.Author]
		if ok {
			continue
		}

		var id int
		err := repo.db.QueryRow(context.Background(),
			"SELECT id FROM users WHERE lower(nickname) = lower($1)", post.Author).Scan(&id)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, domain.ErrNotFound
			}
			return nil, err
		}

		res[post.Author] = id
	}

	return res, nil
}

func (repo *PostRepository) GetPost(id int64) (*domain.Post, error) {
	post := &domain.Post{Id: id}

	err := repo.db.QueryRow(context.Background(),
		`SELECT u.nickname, p.message, p.is_edited,
				f.slug, p.thread_id, p.created_at
		 FROM post p JOIN users u ON u.id = p.author_id
		 			 JOIN thread t ON t.id = p.thread_id
					 JOIN forum f ON f.id = t.forum_id
		 WHERE p.id = $1`, id).
		Scan(
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
		)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return post, nil
}

func (repo *PostRepository) Update(post *domain.Post) error {
	previous, err := repo.GetPost(post.Id)
	if err != nil {
		return err
	}

	if post.Message == "" {
		post.Message = previous.Message
	}

	if post.Message != previous.Message {
		post.IsEdited = true
	}

	_, err = repo.db.Exec(context.Background(),
		`UPDATE post SET message = $1, is_edited = $2 WHERE id = $3`,
		post.Message, post.IsEdited, post.Id)

	if err != nil {
		return err
	}

	post.Author = previous.Author
	post.Forum = previous.Forum
	post.Created = previous.Created
	post.Parent = previous.Parent
	post.Thread = previous.Thread

	return nil
}
