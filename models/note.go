package models

import (
	"context"
	"fmt"
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
	collection     *mongo.Collection
	userCollection *mongo.Collection
}

func NewNoteModel(noteCollection, userCollection *mongo.Collection) *NoteModel {
	return &NoteModel{
		collection:     noteCollection,
		userCollection: userCollection,
	}
}

func (m *NoteModel) Create(userID primitive.ObjectID, title, body string) (*Note, error) {
	// Validate user exists (similar to Order model pattern)
	var userExists struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	err := m.userCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&userExists)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %v", err)
	}

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
		return nil, fmt.Errorf("failed to create note: %v", err)
	}

	note.ID = result.InsertedID.(primitive.ObjectID)
	return note, nil
}

func (m *NoteModel) GetAll(userID primitive.ObjectID) ([]Note, error) {
	var notes []Note

	cursor, err := m.collection.Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		return []Note{}, fmt.Errorf("failed to fetch notes: %v", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var note Note
		if err := cursor.Decode(&note); err != nil {
			return []Note{}, fmt.Errorf("failed to decode note: %v", err)
		}
		notes = append(notes, note)
	}

	if notes == nil {
		return []Note{}, nil
	}

	return notes, nil
}

func (m *NoteModel) GetByID(id primitive.ObjectID, userID primitive.ObjectID) (*Note, error) {
	var note Note
	err := m.collection.FindOne(context.Background(), bson.M{
		"_id":     id,
		"user_id": userID,
	}).Decode(&note)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("note not found")
		}
		return nil, fmt.Errorf("failed to fetch note: %v", err)
	}

	return &note, nil
}

func (m *NoteModel) Update(id primitive.ObjectID, userID primitive.ObjectID, title, body string) (*Note, error) {
	// First check if note exists and belongs to user
	var existingNote Note
	err := m.collection.FindOne(context.Background(), bson.M{
		"_id":     id,
		"user_id": userID,
	}).Decode(&existingNote)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("note not found")
		}
		return nil, fmt.Errorf("failed to find note: %v", err)
	}

	update := bson.M{
		"$set": bson.M{
			"title":      title,
			"body":       body,
			"updated_at": time.Now(),
		},
	}

	result, err := m.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id, "user_id": userID},
		update,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update note: %v", err)
	}

	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no note was updated")
	}

	return m.GetByID(id, userID)
}

func (m *NoteModel) Delete(id primitive.ObjectID, userID primitive.ObjectID) error {
	filter := bson.M{
		"_id":     id,
		"user_id": userID,
	}

	// First check if note exists and belongs to user
	var note Note
	err := m.collection.FindOne(context.Background(), filter).Decode(&note)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("note not found")
		}
		return fmt.Errorf("failed to find note: %v", err)
	}

	result, err := m.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete note: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no note was deleted")
	}

	return nil
}
