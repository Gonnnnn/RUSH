package attendance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongodbAttendance struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	SessionId   string             `bson:"session_id"`
	SessionName string             `bson:"session_name"`
	UserId      string             `bson:"user_id"`
	UserName    string             `bson:"user_name"`
	JoinedAt    time.Time          `bson:"joined_at"`
	CreatedAt   time.Time          `bson:"created_at"`
}

type mongodbRepo struct {
	collection *mongo.Collection
	// The clock to get the current time. It's used to mock the time in tests.
	clock clock.Clock
}

func NewMongoDbRepo(collection *mongo.Collection, clock clock.Clock) *mongodbRepo {
	return &mongodbRepo{
		collection: collection,
		clock:      clock,
	}
}

func (m *mongodbRepo) FindBySessionId(sessionId string) ([]Attendance, error) {
	return nil, errors.New("not implemented")
}

func (m *mongodbRepo) FindByUserId(userId string) ([]Attendance, error) {
	ctx := context.Background()

	cursor, err := m.collection.Find(ctx, bson.M{"user_id": userId})
	if err != nil {
		return nil, fmt.Errorf("failed to query attendances: %w", err)
	}
	defer cursor.Close(ctx)

	var attendances []mongodbAttendance
	if err = cursor.All(ctx, &attendances); err != nil {
		return nil, fmt.Errorf("failed to decode attendances: %w", err)
	}

	var converted []Attendance
	for _, a := range attendances {
		converted = append(converted, Attendance{
			Id:          a.Id.Hex(),
			SessionId:   a.SessionId,
			SessionName: a.SessionName,
			UserId:      a.UserId,
			UserName:    a.UserName,
			JoinedAt:    a.JoinedAt,
			CreatedAt:   a.CreatedAt,
		})
	}

	return converted, nil
}

type AddAttendanceReq struct {
	SessionId   string
	SessionName string
	UserId      string
	UserName    string
	JoinedAt    time.Time
}

func (m *mongodbRepo) BulkInsert(requests []AddAttendanceReq) error {
	if len(requests) == 0 {
		return nil
	}

	attendances := make([]interface{}, 0, len(requests))
	now := m.clock.Now()

	for _, request := range requests {
		attendances = append(attendances, &mongodbAttendance{
			SessionId:   request.SessionId,
			SessionName: request.SessionName,
			UserId:      request.UserId,
			UserName:    request.UserName,
			JoinedAt:    request.JoinedAt,
			CreatedAt:   now,
		})
	}

	_, err := m.collection.InsertMany(context.Background(), attendances)
	return err
}
