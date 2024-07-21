package attendance

import (
	"database/sql"
	"fmt"
	"strings"
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

func (r *repo) GetAll() ([]AttendanceReport, error) {
	rows, err := r.db.Query(`SELECT * FROM attendance_reports`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reports := []AttendanceReport{}
	for rows.Next() {
		var report AttendanceReport
		err := rows.Scan(&report.Id, &report.Name, &report.Description, &report.SessionIds, &report.CreatedAt, &report.CreatedBy)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func (r *repo) Add(name string, description string, sessionIds []string, createdBy int) error {
	_, err := r.db.Exec(
		`INSERT INTO sessions (name, description, session_ids, created_at, created_by) VALUES (?, ?, ?, ?, ?)`,
		name, description, integerList(strings.Join(sessionIds, ",")), time.Now(), createdBy,
	)
	if err != nil {
		return fmt.Errorf("failed to insert session: %w", err)
	}

	return nil
}
