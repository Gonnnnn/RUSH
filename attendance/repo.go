package attendance

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongodbAttendance struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	SessionIds  []string           `bson:"session_ids"`
	CreatedAt   time.Time          `bson:"created_at"`
	CreatedBy   int                `bson:"created_by"`
}

type mongodbRepo struct {
	collection *mongo.Collection
}

type repo struct {
	db *sql.DB
}

func NewMongoDbRepo(collection *mongo.Collection) *mongodbRepo {
	return &mongodbRepo{
		collection: collection,
	}
}

func NewRepo(db *sql.DB) *repo {
	return &repo{
		db: db,
	}
}

func (r *mongodbRepo) GetAll() ([]AttendanceReport, error) {
	ctx := context.Background()
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reports []AttendanceReport
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *mongodbRepo) Add(name string, description string, sessionIds []string, createdBy int) error {
	ctx := context.Background()
	report := AttendanceReport{
		Name:        name,
		Description: description,
		SessionIds:  sessionIds,
		CreatedAt:   time.Now(),
		CreatedBy:   createdBy,
	}

	_, err := r.collection.InsertOne(ctx, report)
	if err != nil {
		return fmt.Errorf("failed to insert attendance report: %w", err)
	}

	return nil
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
		name, description, strings.Join(sessionIds, ","), time.Now(), createdBy,
	)
	if err != nil {
		return fmt.Errorf("failed to insert session: %w", err)
	}

	return nil
}
