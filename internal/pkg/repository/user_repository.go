package repository

import (
	"context"
	"errors"
	"fmt"

	"forum/internal/pkg/domain"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) GetByNickname(nickname string) (*domain.User, error) {
	res := &domain.User{}

	err := repo.db.QueryRow(context.Background(),
		`SELECT id, nickname, fullname, about, email
			FROM users WHERE lower(nickname) = lower($1)`, nickname).
		Scan(
			&res.Id,
			&res.Nickname,
			&res.Fullname,
			&res.About,
			&res.Email,
		)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return res, err
}

func (repo *UserRepository) Create(user *domain.User) (domain.UserBatch, error) {
	_, err := repo.db.Exec(context.Background(),
		`INSERT INTO users(nickname, fullname, about, email)
				VALUES ($1, $2, $3, $4)`,
		user.Nickname, user.Fullname, user.About, user.Email)

	if err == nil {
		return nil, nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != pgerrcode.UniqueViolation {
		return nil, err
	}

	resp := []*domain.User{}

	rows, err := repo.db.Query(context.Background(),
		`SELECT id, nickname, fullname, about, email
		 FROM users WHERE lower(email) = lower($1) OR lower(nickname) = lower($2)`,
		user.Email, user.Nickname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &domain.User{}
		err = rows.Scan(
			&user.Id,
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email,
		)
		if err != nil {
			return nil, err
		}
		resp = append(resp, user)
	}

	return resp, domain.ErrAlreadyExists
}

func (repo *UserRepository) Update(user *domain.User) error {
	previous, err := repo.GetByNickname(user.Nickname)
	if err != nil {
		return err
	}

	if user.Fullname == "" {
		user.Fullname = previous.Fullname
	}
	if user.Email == "" {
		user.Email = previous.Email
	}
	if !user.About.Valid {
		user.About = previous.About
	}

	_, err = repo.db.Exec(context.Background(),
		`UPDATE users SET fullname = $1, about = $2, email = $3 WHERE lower(nickname) = lower($4)`,
		user.Fullname, user.About, user.Email, user.Nickname)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (repo *UserRepository) GetByForum(params *domain.UserListParams) (domain.UserBatch, error) {
	query := `SELECT u.id, u.nickname, u.fullname, u.about, u.email
				FROM users u JOIN user_forum uf ON u.id = uf.user_id
							JOIN forum f ON f.id = uf.forum_id
				WHERE f.id = $1 `

	args := []interface{}{params.ForumId}

	if params.Since != "" {
		query += "AND lower(u.nickname) "
		args = append(args, params.Since)

		if !params.Desc {
			query += "> lower($2)"
		} else {
			query += "< lower($2)"
		}
	}

	query += " ORDER BY lower(u.nickname)"
	if params.Desc {
		query += " DESC"
	}

	args = append(args, params.Limit)
	query += fmt.Sprintf(" LIMIT $%d", len(args))

	rows, err := repo.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}

	users, err := pgx.CollectRows[*domain.User](rows, func(row pgx.CollectableRow) (*domain.User, error) {
		user := &domain.User{}
		err := row.Scan(
			&user.Id,
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email,
		)
		return user, err
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return users, nil
		}
		return nil, err
	}

	return users, nil
}
