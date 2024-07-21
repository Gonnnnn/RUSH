package user

import (
	"database/sql"
)

type repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) GetAll() ([]User, error) {
	rows, err := r.db.Query(`SELECT * FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.University, &user.Phone, &user.Generation, &user.IsActive)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *repo) Add(u *User) error {
	_, err := r.db.Exec(
		`INSERT INTO users (name, university, phone, generation, is_active) VALUES (?, ?, ?, ?, ?)`,
		u.Name, u.University, u.Phone, u.Generation, u.IsActive,
	)
	return err
}
