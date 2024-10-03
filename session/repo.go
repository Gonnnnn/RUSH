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
	Id               primitive.ObjectID `bson:"_id,omitempty"`
	Name             string             `bson:"name"`
	Description      string             `bson:"description"`
	CreatedBy        string             `bson:"created_by"`
	GoogleFormId     string             `bson:"google_form_id"`
	GoogleFormUri    string             `bson:"google_form_uri"`
	CreatedAt        time.Time          `bson:"created_at"`
	StartsAt         time.Time          `bson:"starts_at"`
	Score            int                `bson:"score"`
	AttendanceStatus AttendanceStatus   `bson:"attendance_status"`
	IsDeleted        bool               `bson:"is_deleted"`
}

type mongodbRepo struct {
	collection *mongo.Collection
}

type UpdateForm struct {
	Title         *string
	Description   *string
	GoogleFormId  *string
	GoogleFormUri *string
	// It should be updated with the form's description.
	StartsAt         *time.Time
	Score            *int
	AttendanceStatus *AttendanceStatus

	// Indicator to return the updated session.
	ReturnUpdatedSession bool
}

func NewMongoDbRepo(collection *mongo.Collection) *mongodbRepo {
	return &mongodbRepo{
		collection: collection,
	}
}

func (r *mongodbRepo) Get(id string) (Session, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Session{}, fmt.Errorf("invalid id: %w", err)
	}

	session := &mongodbSession{}
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectID, "is_deleted": false}).Decode(session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Session{}, fmt.Errorf("session not found")
		}
		return Session{}, fmt.Errorf("failed to get session: %w", err)
	}

	return *fromMongodbSession(session), nil
}

func (r *mongodbRepo) GetOpenSessions() ([]Session, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"is_closed": false, "is_deleted": false})
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

func (r *mongodbRepo) GetAll() ([]Session, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"is_deleted": false})
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

	total, err := r.collection.CountDocuments(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, fmt.Errorf("failed to count sessions: %w", err)
	}

	// Fetch pageSize + 1 to check if there are more pages.
	cursor, err := r.collection.Find(ctx, bson.M{"is_deleted": false},
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

func (r *mongodbRepo) Add(name string, description string, createdBy string, startsAt time.Time, score int) (string, error) {
	session := mongodbSession{
		Name:             name,
		Description:      description,
		CreatedBy:        createdBy,
		GoogleFormId:     "",
		GoogleFormUri:    "",
		CreatedAt:        time.Now(),
		StartsAt:         startsAt,
		Score:            score,
		AttendanceStatus: AttendanceStatusNotAppliedYet,
		IsDeleted:        false,
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

func (r *mongodbRepo) Update(id string, updateForm UpdateForm) (Session, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Session{}, fmt.Errorf("invalid id: %w", err)
	}

	update := bson.M{}
	if updateForm.Title != nil {
		update["name"] = *updateForm.Title
	}
	if updateForm.Description != nil {
		update["description"] = *updateForm.Description
	}
	if updateForm.GoogleFormId != nil {
		update["google_form_id"] = *updateForm.GoogleFormId
	}
	if updateForm.GoogleFormUri != nil {
		update["google_form_uri"] = *updateForm.GoogleFormUri
	}
	if updateForm.StartsAt != nil {
		update["starts_at"] = *updateForm.StartsAt
	}
	if updateForm.Score != nil {
		update["score"] = *updateForm.Score
	}
	if updateForm.AttendanceStatus != nil {
		update["attendance_status"] = *updateForm.AttendanceStatus
	}

	if _, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": update}); err != nil {
		return Session{}, fmt.Errorf("failed to update session: %w", err)
	}

	if updateForm.ReturnUpdatedSession {
		return r.Get(id)
	}

	return Session{}, nil
}

func (r *mongodbRepo) Delete(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	_, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": bson.M{"is_deleted": true}})
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func fromMongodbSession(session *mongodbSession) *Session {
	return &Session{
		Id:               session.Id.Hex(),
		Name:             session.Name,
		Description:      session.Description,
		CreatedBy:        session.CreatedBy,
		GoogleFormId:     session.GoogleFormId,
		GoogleFormUri:    session.GoogleFormUri,
		CreatedAt:        session.CreatedAt,
		StartsAt:         session.StartsAt,
		Score:            session.Score,
		AttendanceStatus: session.AttendanceStatus,
	}
}
