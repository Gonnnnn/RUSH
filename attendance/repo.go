package attendance

import (
	"context"
	"errors"
	"time"

	"github.com/benbjohnson/clock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongodbAttendance struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	SessionId string             `bson:"session_id"`
	UserId    string             `bson:"user_id"`
	JoinedAt  time.Time          `bson:"joined_at"`
	CreatedAt time.Time          `bson:"created_at"`
}

type mongodbRepo struct {
	collection *mongo.Collection
	clock      clock.Clock
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
	return nil, errors.New("not implemented")
}

func (m *mongodbRepo) BulkInsert(sessionId string, userIds []string) error {
	if len(userIds) == 0 {
		return nil
	}

	attendances := make([]interface{}, 0, len(userIds))
	now := m.clock.Now()

	for _, userId := range userIds {
		attendance := mongodbAttendance{
			SessionId: sessionId,
			UserId:    userId,
			// TODO(#42): Uses correct values for JoinedAt.
			JoinedAt:  now,
			CreatedAt: now,
		}
		attendances = append(attendances, attendance)
	}

	_, err := m.collection.InsertMany(context.Background(), attendances)
	return err
}
