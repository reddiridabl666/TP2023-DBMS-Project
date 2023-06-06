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

func (repo *UserRepository) GetByNickname(nickname string) (*domain.User, error) {
	res := &domain.User{}

	err := repo.db.QueryRow(
		`SELECT id, nickname, fullname, about, email
			FROM users WHERE lower(nickname) = lower($1)`, nickname).
		Scan(
			&res.Id,
			&res.Nickname,
			&res.Fullname,
			&res.About,
			&res.Email,
		)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return res, err
}

func (repo *UserRepository) Create(user *domain.User) (domain.UserBatch, error) {
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
		 FROM users WHERE lower(email) = lower($1) OR lower(nickname) = lower($2)`,
		user.Email, user.Nickname)
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
			&user.Email,
		)
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

	_, err = repo.db.Exec(
		`UPDATE users SET fullname = $1, about = $2, email = $3 WHERE lower(nickname) = lower($4)`,
		user.Fullname, user.About, user.Email, user.Nickname)
	if err != nil {
		if pgError, ok := err.(pgx.PgError); ok && pgError.Code == pgerrcode.UniqueViolation {
			return domain.ErrAlreadyExists
		}
		return err
	}
	return nil
}

// curl -v -X POST --data '{"about": "Quodam abs en cui. Bene talia ipsum. Subdita. Libenter inludi veni tolerantiam pacto.","email": "alioquin.3ZIhFjfMlviW7@precescit.com","fullname": "Joshua Brown"}' http://localhost:5000/api/user/retarder.lHFMfjC63d5q7u/create
