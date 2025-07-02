package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserRepository interface {
}

type UserRepository struct {
	db *mongo.Collection
}

func NewUserRepository(col *mongo.Collection) IUserRepository {
	return &UserRepository{db: col}
}

func (r *UserRepository) FindUsers()