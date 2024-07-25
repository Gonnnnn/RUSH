package session

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongodbSession struct {
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string             `bson:"name"`
	Description   string             `bson:"description"`
	HostedBy      int                `bson:"hosted_by"`
	CreatedBy     int                `bson:"created_by"`
	JoinningUsers string             `bson:"joinning_users"`
	CreatedAt     time.Time          `bson:"created_at"`
	StartsAt      time.Time          `bson:"starts_at"`
	Score         int                `bson:"score"`
	IsClosed      bool               `bson:"is_closed"`
}

type mongodbRepo struct {
	collection *mongo.Collection
}

type sqliteRepo struct {
	db *sql.DB
}

func NewMongoDBRepo(collection *mongo.Collection) *mongodbRepo {
	return &mongodbRepo{
		collection: collection,
	}
}

func NewSqliteRepo(db *sql.DB) *sqliteRepo {
	return &sqliteRepo{
		db: db,
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

func (r *mongodbRepo) Add(name string, description string, hostedBy int, createdBy int, startsAt time.Time, score int) (string, error) {
	session := mongodbSession{
		Name:          name,
		Description:   description,
		HostedBy:      hostedBy,
		CreatedBy:     createdBy,
		JoinningUsers: "",
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

func (r *sqliteRepo) Get(id string) (*Session, error) {
	session := &Session{}
	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	if r.db.QueryRow(`SELECT * FROM sessions WHERE id = ?`, intId).Scan(
		&session.Id, &session.Name, &session.Description, &session.HostedBy, &session.CreatedBy, &session.JoinningUsers, &session.CreatedAt, &session.StartsAt, &session.Score, &session.IsClosed,
	); err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

func (r *sqliteRepo) GetAll() ([]Session, error) {
	rows, err := r.db.Query(`SELECT * FROM sessions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := []Session{}
	for rows.Next() {
		var session Session
		err := rows.Scan(&session.Id, &session.Name, &session.Description, &session.HostedBy, &session.CreatedBy, &session.JoinningUsers, &session.CreatedAt, &session.StartsAt, &session.Score, &session.IsClosed)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *sqliteRepo) Add(name string, description string, hostedBy int, createdBy int, startsAt time.Time, score int) (string, error) {
	result, err := r.db.Exec(
		`INSERT INTO sessions (name, description, hosted_by, created_by, joinning_users, created_at, starts_at, score, is_closed) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		name, description, hostedBy, createdBy, "", time.Now(), startsAt, score, false,
	)
	if err != nil {
		return "", fmt.Errorf("failed to insert session: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("failed to get last insert id: %w", err)
	}
	return strconv.Itoa(int(id)), nil
}

func fromMongodbSession(session *mongodbSession) *Session {
	return &Session{
		Id:            session.Id.Hex(),
		Name:          session.Name,
		Description:   session.Description,
		HostedBy:      session.HostedBy,
		CreatedBy:     session.CreatedBy,
		JoinningUsers: integerList(session.JoinningUsers),
		CreatedAt:     session.CreatedAt,
		StartsAt:      session.StartsAt,
		Score:         session.Score,
		IsClosed:      session.IsClosed,
	}
}
