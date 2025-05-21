package models

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"passsword" json:"-"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type UserModel struct {
	collection *mongo.Collection
}

func NewUserModel(collection *mongo.Collection) *UserModel {
	return &UserModel{collection: collection}
}

func (m *UserModel) Create(email, password string) (*User, error) {
	// check if user exists
	var existingUser User
	err := m.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&existingUser)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &User{
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	result, err := m.collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil

}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	var user User
	err := m.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) VerifyPassword(user *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
