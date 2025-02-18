package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// The session record in MongoDB.
type mongodbSession struct {
	// The unique identifier for the session. E.g. "1"
	Id primitive.ObjectID `bson:"_id,omitempty"`
	// The name of the session. E.g. "여의도 공원 정규런"
	Name string `bson:"name"`
	// The description of the session. E.g. "여의도 공원에서 정규런을 진행합니다."
	Description string `bson:"description"`
	// The unique identifier for the user who created the session. E.g. "1"
	CreatedBy string `bson:"created_by"`
	// The unique identifier for the Google Form that the session is associated with. E.g. "1"
	GoogleFormId string `bson:"google_form_id"`
	// The URI of the Google Form that the session is associated with. E.g. "https://docs.google.com/forms/d/e/1FAIpQLSf9dFVMN-7HgPXl8jUMyL4ynq-e3fKUZXIQaQ/viewform?usp=sf_link"
	GoogleFormUri string `bson:"google_form_uri"`
	// The time when the session was created. E.g. "2021-01-01T00:00:00Z"
	CreatedAt time.Time `bson:"created_at"`
	// The time when the session starts. E.g. "2021-01-01T00:00:00Z"
	StartsAt time.Time `bson:"starts_at"`
	// The score of the session. E.g. 2
	Score int `bson:"score"`
	// The status of the session's attendance.
	// It indicates if it is applied, ignored, etc.
	AttendanceStatus AttendanceStatus `bson:"attendance_status"`
	// The reason why the attendance is ignored. E.g. "The user is not a member."
	AttendanceIgnoredReason string `bson:"attendance_ignored_reason"`
	// Whether the session is deleted. E.g. false
	IsDeleted bool `bson:"is_deleted"`
}

type mongodbRepo struct {
	collection *mongo.Collection
}

var ErrNotFound = errors.New("session not found")

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
			return Session{}, ErrNotFound
		}
		return Session{}, fmt.Errorf("failed to get session: %w", err)
	}

	return *fromMongodbSession(session), nil
}

// Get open sessions that has its attendance form. Open means the session has not closed, as in the attendance
// is not applied yet.
func (r *mongodbRepo) GetOpenSessionsWithForm() ([]Session, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"attendance_status": "not_applied_yet", "is_deleted": false, "google_form_id": bson.M{"$ne": ""}})
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

// List sessions with pagination.
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

// The form to update the session. It only includes fields that can be updated.
type UpdateForm struct {
	Title                   *string
	Description             *string
	GoogleFormId            *string
	GoogleFormUri           *string
	StartsAt                *time.Time
	Score                   *int
	AttendanceStatus        *AttendanceStatus
	AttendanceIgnoredReason *string

	// Indicator to return the updated session. If false, the session is not returned.
	ReturnUpdatedSession bool
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
	if updateForm.AttendanceIgnoredReason != nil {
		update["attendance_ignored_reason"] = *updateForm.AttendanceIgnoredReason
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
