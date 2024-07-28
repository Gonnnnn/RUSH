package user

import (
	"context"
	"database/sql"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongodbUser struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	University string             `bson:"university"`
	Phone      string             `bson:"phone"`
	Generation string             `bson:"generation"`
	IsActive   bool               `bson:"is_active"`
}

type mongodbRepo struct {
	collection *mongo.Collection
}

type repo struct {
	db *sql.DB
}

func NewMongoDbRepo(collection *mongo.Collection) *mongodbRepo {
	return &mongodbRepo{
		collection: collection,
	}
}

func NewRepo(db *sql.DB) *repo {
	return &repo{
		db: db,
	}
}

func (r *mongodbRepo) GetAll() ([]User, error) {
	ctx := context.Background()
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []mongodbUser
	if err = cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	var converted []User
	for _, u := range users {
		converted = append(converted, User{
			Id:         u.Id.Hex(),
			Name:       u.Name,
			University: u.University,
			Phone:      u.Phone,
			Generation: u.Generation,
			IsActive:   u.IsActive,
		})
	}

	return converted, nil
}

func (r *mongodbRepo) Add(u *User) error {
	ctx := context.Background()
	_, err := r.collection.InsertOne(ctx, u)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (r *repo) GetAll() ([]User, error) {
	rows, err := r.db.Query(`SELECT * FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.University, &user.Phone, &user.Generation, &user.IsActive)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *repo) Add(u *User) error {
	_, err := r.db.Exec(
		`INSERT INTO users (name, university, phone, generation, is_active) VALUES (?, ?, ?, ?, ?)`,
		u.Name, u.University, u.Phone, u.Generation, u.IsActive,
	)
	return err
}
