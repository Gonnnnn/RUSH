package session

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) Get(id string) (*Session, error) {
	session := &Session{}
	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	if r.db.QueryRow(`SELECT * FROM sessions WHERE id = ?`, intId).Scan(
		&session.Id, &session.Name, &session.Description, &session.HostedBy, &session.CreatedBy, &session.JoinningUsers, &session.CreatedAt, &session.StartsAt, &session.Score, &session.IsClosed,
	); err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

func (r *repo) GetAll() ([]Session, error) {
	rows, err := r.db.Query(`SELECT * FROM sessions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := []Session{}
	for rows.Next() {
		var session Session
		err := rows.Scan(&session.Id, &session.Name, &session.Description, &session.HostedBy, &session.CreatedBy, &session.JoinningUsers, &session.CreatedAt, &session.StartsAt, &session.Score, &session.IsClosed)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *repo) Add(name string, description string, hostedBy int, createdBy int, startsAt time.Time, score int) (string, error) {
	result, err := r.db.Exec(
		`INSERT INTO sessions (name, description, hosted_by, created_by, joinning_users, created_at, starts_at, score, is_closed) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		name, description, hostedBy, createdBy, "", time.Now(), startsAt, score, false,
	)
	if err != nil {
		return "", fmt.Errorf("failed to insert session: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("failed to get last insert id: %w", err)
	}
	return strconv.Itoa(int(id)), nil
}
