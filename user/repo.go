package user

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodbUser struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	University   string             `bson:"university"`
	Phone        string             `bson:"phone"`
	Generation   float64            `bson:"generation"`
	IsActive     bool               `bson:"is_active"`
	ExternalName string             `bson:"external_name"`
}

type mongodbRepo struct {
	collection *mongo.Collection
}

var ErrNotFound = errors.New("user not found")

func NewMongoDbRepo(collection *mongo.Collection) *mongodbRepo {
	return &mongodbRepo{
		collection: collection,
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
			Id:           u.Id.Hex(),
			Name:         u.Name,
			University:   u.University,
			Phone:        u.Phone,
			Generation:   u.Generation,
			IsActive:     u.IsActive,
			ExternalName: u.ExternalName,
		})
	}

	return converted, nil
}

func (r *mongodbRepo) GetByEmail(email string) (*User, error) {
	ctx := context.Background()

	var u mongodbUser
	if err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &User{
		Id:           u.Id.Hex(),
		Name:         u.Name,
		University:   u.University,
		Phone:        u.Phone,
		Generation:   u.Generation,
		IsActive:     u.IsActive,
		ExternalName: u.ExternalName,
	}, nil
}

type ListResult struct {
	Users      []User `json:"users"`
	IsEnd      bool   `json:"is_end"`
	TotalCount int    `json:"total_count"`
}

func (r *mongodbRepo) List(offset int, pageSize int) (*ListResult, error) {
	ctx := context.Background()

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Fetch pageSize + 1 to check if there are more pages.
	cursor, err := r.collection.Find(ctx, bson.M{},
		options.Find().
			SetSkip(int64(offset)).
			SetLimit(int64(pageSize+1)).
			SetSort(bson.D{{Key: "generation", Value: -1}}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []mongodbUser
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	isEnd := len(users) <= pageSize
	if !isEnd {
		users = users[:pageSize]
	}

	converted := make([]User, len(users))
	for i, u := range users {
		converted[i] = User{
			Id:           u.Id.Hex(),
			Name:         u.Name,
			University:   u.University,
			Phone:        u.Phone,
			Generation:   u.Generation,
			IsActive:     u.IsActive,
			ExternalName: u.ExternalName,
		}
	}

	return &ListResult{
		Users:      converted,
		IsEnd:      isEnd,
		TotalCount: int(total),
	}, nil
}

func (r *mongodbRepo) Get(id string) (*User, error) {
	ctx := context.Background()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to convert id: %w", err)
	}

	var u mongodbUser
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &User{
		Id:           u.Id.Hex(),
		Name:         u.Name,
		University:   u.University,
		Phone:        u.Phone,
		Generation:   u.Generation,
		IsActive:     u.IsActive,
		ExternalName: u.ExternalName,
	}, nil
}

func (r *mongodbRepo) CountByName(name string) (int, error) {
	ctx := context.Background()
	count, err := r.collection.CountDocuments(ctx, bson.M{"name": name})
	if err != nil {
		return 0, fmt.Errorf("database client has failed: %w", err)
	}
	return int(count), nil
}

func (r *mongodbRepo) GetAllByExternalNames(externalNames []string) ([]User, error) {
	ctx := context.Background()

	var users []mongodbUser
	cursor, err := r.collection.Find(ctx, bson.M{"external_name": bson.M{"$in": externalNames}})
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	converted := make([]User, len(users))
	for index, user := range users {
		converted[index] = User{
			Id:           user.Id.Hex(),
			Name:         user.Name,
			University:   user.University,
			Phone:        user.Phone,
			Generation:   user.Generation,
			IsActive:     user.IsActive,
			ExternalName: user.ExternalName,
		}
	}

	return converted, nil
}

func (r *mongodbRepo) Add(user User) error {
	ctx := context.Background()

	_, err := r.collection.InsertOne(ctx, mongodbUser{
		Name:         user.Name,
		University:   user.University,
		Phone:        user.Phone,
		Generation:   user.Generation,
		IsActive:     user.IsActive,
		ExternalName: user.ExternalName,
	})
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}
