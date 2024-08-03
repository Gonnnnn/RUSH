package attendance

import (
	"errors"
	"time"

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
}

func NewMongoDbRepo(collection *mongo.Collection) *mongodbRepo {
	return &mongodbRepo{
		collection: collection,
	}
}

func FindBySessionId(sessionId string) ([]Attendance, error) {
	return nil, errors.New("not implemented")
}

func FindByUserId(userId string) ([]Attendance, error) {
	return nil, errors.New("not implemented")
}

func BulkInsert(sessionId string, userIds []string) error {
	return errors.New("not implemented")
}
