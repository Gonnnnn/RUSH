package session

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodbSession struct {
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string             `bson:"name"`
	Description   string             `bson:"description"`
	HostedBy      int                `bson:"hosted_by"`
	CreatedBy     int                `bson:"created_by"`
	GoogleFormId  string             `bson:"google_form_id"`
	GoogleFormUri string             `bson:"google_form_uri"`
	JoinningUsers []string           `bson:"joinning_users"`
	CreatedAt     time.Time          `bson:"created_at"`
	StartsAt      time.Time          `bson:"starts_at"`
	Score         int                `bson:"score"`
	IsClosed      bool               `bson:"is_closed"`
}

type mongodbRepo struct {
	collection *mongo.Collection
}

type UpdateForm struct {
	Title         *string
	Description   *string
	HostedBy      *int
	GoogleFormId  *string
	GoogleFormUri *string
	JoinningUsers *string
	// It should be updated with the form's description.
	StartsAt *time.Time
	Score    *int
	IsClosed *bool

	// Indicator to return the updated session.
	ReturnUpdatedSession bool
}

func NewMongoDbRepo(collection *mongo.Collection) *mongodbRepo {
	return &mongodbRepo{
		collection: collection,
	}
}

func (r *mongodbRepo) Get(id string) (*Session, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	session := &mongodbSession{}
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return fromMongodbSession(session), nil
}

func (r *mongodbRepo) GetAll() ([]Session, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}
	defer cursor.Close(context.Background())

	var mongoSessions []mongodbSession
	if err = cursor.All(context.Background(), &mongoSessions); err != nil {
		return nil, fmt.Errorf("failed to decode sessions: %w", err)
	}

	sessions := []Session{}
	for _, mongoSession := range mongoSessions {
		sessions = append(sessions, *fromMongodbSession(&mongoSession))
	}
	return sessions, nil
}

type ListResult struct {
	Sessions   []Session
	IsEnd      bool
	TotalCount int
}

func (r *mongodbRepo) List(offset int, pageSize int) (*ListResult, error) {
	ctx := context.Background()

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to count sessions: %w", err)
	}

	// Fetch pageSize + 1 to check if there are more pages.
	cursor, err := r.collection.Find(ctx, bson.M{},
		options.Find().
			SetSkip(int64(offset)).
			SetLimit(int64(pageSize+1)).
			SetSort(bson.D{{Key: "starts_at", Value: -1}}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find sessions: %w", err)
	}
	defer cursor.Close(ctx)

	var mongodbSessions []mongodbSession
	if err := cursor.All(ctx, &mongodbSessions); err != nil {
		return nil, fmt.Errorf("failed to decode sessions: %w", err)
	}

	isEnd := len(mongodbSessions) <= pageSize
	if !isEnd {
		mongodbSessions = mongodbSessions[:pageSize]
	}

	sessions := []Session{}
	for _, mongoSession := range mongodbSessions {
		sessions = append(sessions, *fromMongodbSession(&mongoSession))
	}

	return &ListResult{
		Sessions:   sessions,
		IsEnd:      isEnd,
		TotalCount: int(total),
	}, nil
}

func (r *mongodbRepo) Add(name string, description string, hostedBy int, createdBy int, startsAt time.Time, score int) (string, error) {
	session := mongodbSession{
		Name:          name,
		Description:   description,
		HostedBy:      hostedBy,
		CreatedBy:     createdBy,
		GoogleFormId:  "",
		GoogleFormUri: "",
		JoinningUsers: []string{},
		CreatedAt:     time.Now(),
		StartsAt:      startsAt,
		Score:         score,
		IsClosed:      false,
	}

	result, err := r.collection.InsertOne(context.Background(), session)
	if err != nil {
		return "", fmt.Errorf("failed to insert session: %w", err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to get inserted id")
	}

	return id.Hex(), nil
}

func (r *mongodbRepo) Update(id string, updateForm *UpdateForm) (*Session, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	update := bson.M{}
	if updateForm.Title != nil {
		update["name"] = *updateForm.Title
	}
	if updateForm.Description != nil {
		update["description"] = *updateForm.Description
	}
	if updateForm.HostedBy != nil {
		update["hosted_by"] = *updateForm.HostedBy
	}
	if updateForm.GoogleFormId != nil {
		update["google_form_id"] = *updateForm.GoogleFormId
	}
	if updateForm.GoogleFormUri != nil {
		update["google_form_uri"] = *updateForm.GoogleFormUri
	}
	if updateForm.JoinningUsers != nil {
		update["joinning_users"] = *updateForm.JoinningUsers
	}
	if updateForm.StartsAt != nil {
		update["starts_at"] = *updateForm.StartsAt
	}
	if updateForm.Score != nil {
		update["score"] = *updateForm.Score
	}
	if updateForm.IsClosed != nil {
		update["is_closed"] = *updateForm.IsClosed
	}

	if _, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": update}); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	if updateForm.ReturnUpdatedSession {
		return r.Get(id)
	}

	return nil, nil
}

func fromMongodbSession(session *mongodbSession) *Session {
	return &Session{
		Id:            session.Id.Hex(),
		Name:          session.Name,
		Description:   session.Description,
		HostedBy:      session.HostedBy,
		CreatedBy:     session.CreatedBy,
		GoogleFormId:  session.GoogleFormId,
		GoogleFormUri: session.GoogleFormUri,
		JoinningUsers: session.JoinningUsers,
		CreatedAt:     session.CreatedAt,
		StartsAt:      session.StartsAt,
		Score:         session.Score,
		IsClosed:      session.IsClosed,
	}
}
