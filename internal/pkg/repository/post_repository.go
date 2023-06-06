package repository

import (
	"database/sql"
	"fmt"
	"time"

	"forum/internal/pkg/domain"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (repo *PostRepository) AddPosts(thread *domain.Thread, posts domain.PostBatch) error {
	fmt.Println("___POST ADD BEGIN")
	defer fmt.Println("___POST ADD END")

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

	rows, err := repo.db.Query(query, args...)
	fmt.Println("________________ERROR AFTER POST ADD", err)

	if err != nil {
		pgError, ok := err.(pgx.PgError)

		if !ok {
			return err
		}

		switch pgError.Code {
		case pgerrcode.ForeignKeyViolation:
			return domain.ErrNoParent
		case pgerrcode.IntegrityConstraintViolation:
			return domain.ErrInvalidParent
		}
		return err
	}

	defer rows.Close()

	i = 0
	for rows.Next() {
		err := rows.Scan(&posts[i].Id)
		if err != nil {
			fmt.Println("________________ERROR IN SCAN", err)
			return err
		}
		i++
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
		err := repo.db.QueryRow("SELECT id FROM users WHERE lower(nickname) = lower($1)", post.Author).Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, domain.ErrNotFound
			}
			return nil, err
		}

		res[post.Author] = id
	}

	return res, nil
}
