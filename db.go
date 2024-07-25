package main

import (
	"database/sql"
)

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		university TEXT,
		phone TEXT,
		generation TEXT,
		is_active BOOLEAN
	);

	CREATE TABLE IF NOT EXISTS attendance_reports (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		description TEXT,
		session_ids TEXT,
		create_at TIMESTAMP,
		created_by INTEGER
	);
	`)
	return err
}

func createDummyData(db *sql.DB) error {
	// Insert dummy data only if all the tables are empty
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		_, err := db.Exec(`
		INSERT INTO users (name, university, phone, generation, is_active) VALUES ('Geon Kim', 'Yonsei', '1234567890', '2020', 1);
		`)
		if err != nil {
			return err
		}
	}

	if err := db.QueryRow("SELECT COUNT(*) FROM attendance_reports").Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		_, err := db.Exec(`
		INSERT INTO attendance_reports (name, description, session_ids, create_at, created_by) VALUES ('Report 1', 'Description 1', '1', '2021-01-01 00:00:00', 1);
		`)
		if err != nil {
			return err
		}
	}

	return nil
}
