package aggregation

import (
	"context"
	"errors"
	"time"

	"github.com/benbjohnson/clock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongodbAggregation struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	SessionIds []string           `bson:"session_ids"`
	UserScores []UserScorePair    `bson:"user_scores"`
	CreatedAt  time.Time          `bson:"created_at"`
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

func (m *mongodbRepo) AddAggregation(sessionIds []string, userScores []UserScorePair) (Aggregation, error) {
	now := m.clock.Now()
	aggregation := mongodbAggregation{
		SessionIds: sessionIds,
		UserScores: userScores,
		CreatedAt:  now,
	}

	result, err := m.collection.InsertOne(context.Background(), aggregation)
	if err != nil {
		return Aggregation{}, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return Aggregation{}, errors.New("failed to get inserted ID")
	}

	return Aggregation{
		Id:         id.Hex(),
		SessionIds: sessionIds,
		UserScores: userScores,
		CreatedAt:  now,
	}, nil
}
