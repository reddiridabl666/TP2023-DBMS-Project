package repository

import (
	"database/sql"

	"forum/internal/pkg/domain"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) GetUser(nickname string) (*domain.User, error) {
	res := &domain.User{Nickname: nickname}

	err := repo.db.QueryRow(
		`SELECT id, fullname, about, email
			FROM users WHERE lower(nickname) = lower($1)`, nickname).
		Scan(
			&res.Id,
			&res.Fullname,
			&res.About,
			&res.Email,
		)

	return res, err
}

func (repo *UserRepository) CreateUser(user *domain.User) ([]*domain.User, error) {
	_, err := repo.db.Exec(
		`INSERT INTO users(nickname, fullname, about, email)
				VALUES ($1, $2, $3, $4)`,
		user.Nickname, user.Fullname, user.About, user.Email)

	if err == nil {
		return nil, nil
	}

	pgError, ok := err.(pgx.PgError)
	if !ok || pgError.Code != pgerrcode.UniqueViolation {
		return nil, err
	}

	resp := []*domain.User{}

	rows, err := repo.db.Query(
		`SELECT id, nickname, fullname, about, email
		 FROM users WHERE email = $1 OR lower(name) = lower($2)`,
		user.Email, user.Nickname)
	if err == sql.ErrNoRows {
		return resp, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &domain.User{}
		rows.Scan(
			&user.Id,
			&user.Nickname,
			&user.Fullname,
			&user.About,
		)
		resp = append(resp, user)
	}

	return resp, nil
}

func (repo *UserRepository) UpdateUser(user *domain.User) error {
	res, err := repo.db.Exec(
		`UPDATE users SET fullname = $1, about = $2, email = $3 WHERE lower(nickname) = lower($4)`,
		user.Fullname, user.About, user.Email, user.Nickname)
	if err != nil {
		if pgError, ok := err.(pgx.PgError); ok && pgError.Code == pgerrcode.UniqueViolation {
			return domain.ErrUniqueViolation
		}
		return err
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowCount == 0 {
		return domain.ErrNotFound
	}
	return nil
}
