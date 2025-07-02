package repository

import (
	"context"

	"github.com/bekbek22/JaiYenMarket_backend/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IAuthRepository interface {
	CreateUser(user *model.User) error
	FindByEmail(email string) (*model.User, error)
}

type AuthRepository struct {
	db *mongo.Collection
}

func NewAuthRepository(col *mongo.Collection) IAuthRepository {
	return &AuthRepository{db: col}
}

func (r *AuthRepository) CreateUser(user *model.User) error {
	_, err := r.db.InsertOne(context.TODO(), user)
	return err
}

func (r *AuthRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
