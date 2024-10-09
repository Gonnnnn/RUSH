package attendance

import (
	"context"
	"fmt"
	"rush/golang/array"
	"time"

	"github.com/benbjohnson/clock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodbAttendance struct {
	Id               primitive.ObjectID `bson:"_id,omitempty"`
	SessionId        string             `bson:"session_id"`
	SessionName      string             `bson:"session_name"`
	SessionScore     int                `bson:"session_score"`
	SessionStartedAt time.Time          `bson:"session_started_at"`
	UserId           string             `bson:"user_id"`
	UserExternalName string             `bson:"user_external_name"`
	UserJoinedAt     time.Time          `bson:"user_joined_at"`
	UserGeneration   float64            `bson:"user_generation"`
	CreatedAt        time.Time          `bson:"created_at"`
	CreatedBy        string             `bson:"created_by"`
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

func (m *mongodbRepo) GetAll() ([]Attendance, error) {
	ctx := context.Background()

	cursor, err := m.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to query attendances: %w", err)
	}

	var attendances []mongodbAttendance
	if err = cursor.All(ctx, &attendances); err != nil {
		return nil, fmt.Errorf("failed to decode attendances: %w", err)
	}

	return array.Map(attendances, toAttendance), nil
}

func (m *mongodbRepo) FindBySessionId(sessionId string) ([]Attendance, error) {
	ctx := context.Background()

	cursor, err := m.collection.Find(
		ctx, bson.M{"session_id": sessionId},
		options.Find().SetSort(bson.D{{Key: "user_joined_at", Value: 1}}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query attendances: %w", err)
	}
	defer cursor.Close(ctx)

	var attendances []mongodbAttendance
	if err = cursor.All(ctx, &attendances); err != nil {
		return nil, fmt.Errorf("failed to decode attendances: %w", err)
	}

	return array.Map(attendances, toAttendance), nil
}

func (m *mongodbRepo) FindByUserId(userId string) ([]Attendance, error) {
	ctx := context.Background()

	cursor, err := m.collection.Find(
		ctx, bson.M{"user_id": userId},
		options.Find().SetSort(bson.D{{Key: "session_started_at", Value: -1}}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query attendances: %w", err)
	}
	defer cursor.Close(ctx)

	var attendances []mongodbAttendance
	if err = cursor.All(ctx, &attendances); err != nil {
		return nil, fmt.Errorf("failed to decode attendances: %w", err)
	}

	return array.Map(attendances, toAttendance), nil
}

type AddAttendanceReq struct {
	SessionId        string
	SessionName      string
	SessionScore     int
	SessionStartedAt time.Time
	UserId           string
	UserExternalName string
	UserGeneration   float64
	UserJoinedAt     time.Time
	CreatedBy        string
}

func (m *mongodbRepo) BulkInsert(requests []AddAttendanceReq) error {
	if len(requests) == 0 {
		return nil
	}

	attendances := make([]interface{}, 0, len(requests))
	now := m.clock.Now()

	for _, request := range requests {
		attendances = append(attendances, &mongodbAttendance{
			SessionId:        request.SessionId,
			SessionName:      request.SessionName,
			SessionScore:     request.SessionScore,
			SessionStartedAt: request.SessionStartedAt,
			UserId:           request.UserId,
			UserExternalName: request.UserExternalName,
			UserGeneration:   request.UserGeneration,
			UserJoinedAt:     request.UserJoinedAt,
			CreatedAt:        now,
			CreatedBy:        request.CreatedBy,
		})
	}

	_, err := m.collection.InsertMany(context.Background(), attendances)
	return err
}

func toAttendance(attendance mongodbAttendance) Attendance {
	return Attendance{
		Id:               attendance.Id.Hex(),
		SessionId:        attendance.SessionId,
		SessionName:      attendance.SessionName,
		SessionScore:     attendance.SessionScore,
		SessionStartedAt: attendance.SessionStartedAt,
		UserId:           attendance.UserId,
		UserExternalName: attendance.UserExternalName,
		UserGeneration:   attendance.UserGeneration,
		UserJoinedAt:     attendance.UserJoinedAt,
		CreatedAt:        attendance.CreatedAt,
		CreatedBy:        attendance.CreatedBy,
	}
}
