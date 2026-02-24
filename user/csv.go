package user

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"rush/permission"
	"strconv"
	"strings"
)

// ParseCSV reads a CSV file with columns (name, external_name, generation, email)
// and returns a slice of User structs with role=member and is_active=true.
func ParseCSV(path string) ([]User, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	// Skip header row.
	if _, err := r.Read(); err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}

	var users []User
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read row: %w", err)
		}
		if len(record) != 4 {
			return nil, fmt.Errorf("expected 4 columns, got %d: %v", len(record), record)
		}

		name := strings.TrimSpace(record[0])
		externalName := strings.TrimSpace(record[1])
		genStr := strings.TrimSpace(record[2])
		email := strings.TrimSpace(record[3])

		generation, err := strconv.ParseFloat(genStr, 64)
		if err != nil {
			return nil, fmt.Errorf("parse generation %q for %s: %w", genStr, name, err)
		}

		users = append(users, User{
			Name:         name,
			ExternalName: externalName,
			Generation:   generation,
			Email:        email,
			Role:         permission.RoleMember,
			IsActive:     true,
		})
	}

	return users, nil
}
