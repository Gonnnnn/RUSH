package user

import (
	"context"
	"errors"
	"fmt"
	"rush/permission"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongodbUser is the user model for MongoDB.
type mongodbUser struct {
	// The unique identifier of the user.
	Id primitive.ObjectID `bson:"_id,omitempty"`
	// The name of the user. E.g., "김건"
	Name string `bson:"name"`
	// The role of the user. E.g., "member"
	Role string `bson:"role"`
	// The generation of the user. E.g., 9
	Generation float64 `bson:"generation"`
	// The activity status of the user. E.g., true
	IsActive bool `bson:"is_active"`
	// The email address of the user. E.g., "kim.geon@gmail.com"
	Email string `bson:"email"`
	// The unique name consisting of the user name and a number.
	// It's used as an external ID for the users so that
	// it's easier for them to identify themselves such as in Google Forms.
	ExternalName string `bson:"external_name"`
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
	for _, user := range users {
		convertedUser, err := convertToUser(user)
		if err != nil {
			return nil, fmt.Errorf("failed to convert user: %w", err)
		}
		converted = append(converted, convertedUser)
	}

	return converted, nil
}

func (r *mongodbRepo) GetAllActive() ([]User, error) {
	ctx := context.Background()
	cursor, err := r.collection.Find(ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []mongodbUser
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	converted := make([]User, len(users))
	for index, user := range users {
		convertedUser, err := convertToUser(user)
		if err != nil {
			return nil, fmt.Errorf("failed to convert user: %w", err)
		}
		converted[index] = convertedUser
	}

	return converted, nil
}

// Returns the user by the given email.
// If not found, it returns ErrNotFound.
func (r *mongodbRepo) GetByEmail(email string) (*User, error) {
	ctx := context.Background()

	var user mongodbUser
	if err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	convertedUser, err := convertToUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to convert user: %w", err)
	}
	return &convertedUser, nil
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
	for index, user := range users {
		convertedUser, err := convertToUser(user)
		if err != nil {
			return nil, fmt.Errorf("failed to convert user: %w", err)
		}
		converted[index] = convertedUser
	}

	return &ListResult{
		Users:      converted,
		IsEnd:      isEnd,
		TotalCount: int(total),
	}, nil
}

// Returns the user by the given ID.
// If not found, it returns ErrNotFound.
func (r *mongodbRepo) Get(id string) (*User, error) {
	ctx := context.Background()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to convert id: %w", err)
	}

	var user mongodbUser
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	convertedUser, err := convertToUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to convert user: %w", err)
	}
	return &convertedUser, nil
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
		convertedUser, err := convertToUser(user)
		if err != nil {
			return nil, fmt.Errorf("failed to convert user: %w", err)
		}
		converted[index] = convertedUser
	}

	return converted, nil
}

func (r *mongodbRepo) Add(user User) error {
	ctx := context.Background()

	_, err := r.collection.InsertOne(ctx, mongodbUser{
		Name:         user.Name,
		Generation:   user.Generation,
		IsActive:     user.IsActive,
		Email:        user.Email,
		ExternalName: user.ExternalName,
	})
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

// UpdateForm is the form to update the user.
type UpdateForm struct {
	Name         *string
	Role         *string
	Generation   *float64
	IsActive     *bool
	Email        *string
	ExternalName *string
}

func (r *mongodbRepo) Update(id string, updateForm UpdateForm) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	update := bson.M{}

	if updateForm.Name != nil {
		update["name"] = *updateForm.Name
	}

	if updateForm.Role != nil {
		update["role"] = *updateForm.Role
	}

	if updateForm.Generation != nil {
		update["generation"] = *updateForm.Generation
	}

	if updateForm.IsActive != nil {
		update["is_active"] = *updateForm.IsActive
	}

	if updateForm.Email != nil {
		update["email"] = *updateForm.Email
	}

	if updateForm.ExternalName != nil {
		update["external_name"] = *updateForm.ExternalName
	}

	if _, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": update}); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func convertToUser(user mongodbUser) (User, error) {
	userRole, err := convertRole(user.Role)
	if err != nil {
		return User{}, err
	}
	return User{
		Id:           user.Id.Hex(),
		Name:         user.Name,
		Role:         userRole,
		Generation:   user.Generation,
		IsActive:     user.IsActive,
		Email:        user.Email,
		ExternalName: user.ExternalName,
	}, nil
}

func convertRole(role string) (permission.Role, error) {
	switch permission.Role(role) {
	case permission.RoleSuperAdmin:
		return permission.RoleSuperAdmin, nil
	case permission.RoleAdmin:
		return permission.RoleAdmin, nil
	case permission.RoleMember:
		return permission.RoleMember, nil
	default:
		return "", fmt.Errorf("invalid role: %s", role)
	}
}
