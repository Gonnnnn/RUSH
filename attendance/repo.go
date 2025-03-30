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

// The attendance record in MongoDB.
type mongodbAttendance struct {
	// The unique identifier for the attendance record. E.g. "1"
	Id primitive.ObjectID `bson:"_id,omitempty"`
	// The unique identifier for the session that the user joined. E.g. "1"
	SessionId string `bson:"session_id"`
	// The name of the session that the user joined. E.g. "여의도 공원 정규런"
	SessionName string `bson:"session_name"`
	// The score of the session. E.g. 2
	SessionScore int `bson:"session_score"`
	// The time when the session started.
	SessionStartedAt time.Time `bson:"session_started_at"`
	// The unique identifier for the user. E.g. "1"
	UserId string `bson:"user_id"`
	// The external name of the user. E.g. "김건3"
	UserExternalName string `bson:"user_external_name"`
	// The generation of the user. E.g. 9.5
	UserGeneration float64 `bson:"user_generation"`
	// The time when the user joined the session.
	UserJoinedAt time.Time `bson:"user_joined_at"`
	// The time when the attendance record was created.
	CreatedAt time.Time `bson:"created_at"`
	// The user or service that created the attendance record.
	// E.g. "auto-syncer", "user-id-123"
	CreatedBy string `bson:"created_by"`
	// Whether the attendance record was force applied.
	ForceApply bool `bson:"force_apply"`
}

type mongodbRepo struct {
	// The actual client that executes the queries.
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

// The request to add an attendance record.
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
	ForceApply       bool
}

func (m *mongodbRepo) BulkInsert(requests []AddAttendanceReq) error {
	if len(requests) == 0 {
		return nil
	}

	// interface type because InsertMany requires []interface{}.
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
			ForceApply:       request.ForceApply,
		})
	}

	_, err := m.collection.InsertMany(context.Background(), attendances)
	return err
}

// The form to update the attendance record of a user.
type UpdateUserAttendanceForm struct {
	// New external name of the user. Leave it nil to not update.
	UserExternalName *string
	// New generation of the user. Leave it nil to not update.
	UserGeneration *float64
}

// Update the information about the user through all of the attendance records of the user.
func (m *mongodbRepo) UpdateUserAttendance(userId string, updateForm UpdateUserAttendanceForm) error {
	update := bson.M{}
	if updateForm.UserExternalName != nil {
		update["user_external_name"] = *updateForm.UserExternalName
	}
	if updateForm.UserGeneration != nil {
		update["user_generation"] = *updateForm.UserGeneration
	}

	_, err := m.collection.UpdateMany(context.Background(), bson.M{"user_id": userId}, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update attendances: %w", err)
	}

	return nil
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
