package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Note struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title     string             `bson:"title" json:"title"`
	Body      string             `bson:"body" json:"body"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type NoteModel struct {
	collection *mongo.Collection
}

func NewNoteModel(collection *mongo.Collection) *NoteModel {
	return &NoteModel{collection: collection}
}

func (m *NoteModel) Create(userID primitive.ObjectID, title, body string) (*Note, error) {
	now := time.Now()
	note := &Note{
		UserID:    userID,
		Title:     title,
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := m.collection.InsertOne(context.Background(), note)
	if err != nil {
		return nil, err
	}

	note.ID = result.InsertedID.(primitive.ObjectID)
	return note, nil
}

func (m *NoteModel) GetAll(userID primitive.ObjectID) ([]Note, error) {
	cursor, err := m.collection.Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var notes []Note
	if err = cursor.All(context.Background(), &notes); err != nil {
		return nil, err
	}

	return notes, nil
}

func (m *NoteModel) GetByID(id primitive.ObjectID, userID primitive.ObjectID) (*Note, error) {
	var note Note
	err := m.collection.FindOne(context.Background(), bson.M{"_id": id, "user_id": userID}).Decode(&note)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (m *NoteModel) Update(id primitive.ObjectID, userID primitive.ObjectID, title, body string) (*Note, error) {
	update := bson.M{
		"$set": bson.M{
			"title":      title,
			"body":       body,
			"updated_at": time.Now(),
		},
	}

	_, err := m.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id, "user_id": userID},
		update,
	)
	if err != nil {
		return nil, err
	}

	return m.GetByID(id, userID)
}

func (m *NoteModel) Delete(id primitive.ObjectID, userID primitive.ObjectID) error {
	_, err := m.collection.DeleteOne(context.Background(), bson.M{"_id": id, "user_id": userID})
	return err
}
