package sendgrid

import (
	"context"
	"log"

	"github.com/jayco/go-local-email/internal/store"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

const sgCollection = "sendgrid"

// SGCollection ..
type SGCollection struct {
	Mail *mongo.Collection
}

// Sendgrid client
type Sendgrid interface {
	Set(ctx context.Context, item *SGItem)
	GetByID(ctx context.Context, id *string) *SGItem
	GetByEmail(ctx context.Context, email *string) *SGItem
	Query(ctx context.Context, filter bson.M) *SGItem
	All(ctx context.Context) []SGItem
	Delete(ctx context.Context, item *SGItem) error
}

// Set an item in the store
func (s *SGCollection) Set(ctx context.Context, item *SGItem) {
	_, err := s.Mail.InsertOne(ctx, item)
	if err != nil {
		log.Println(err)
	}
}

// Query the collection to return a result from a given filter
func (s *SGCollection) Query(ctx context.Context, filter bson.M) *SGItem {
	var item SGItem

	err := s.Mail.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &item
}

// GetByID returns a stored item by a given id or nil if it can not be found
func (s *SGCollection) GetByID(ctx context.Context, id *string) *SGItem {
	return s.Query(ctx, bson.M{"_id": *id})
}

// GetByEmail returns a stored item by email address or nil if it can not be found
func (s *SGCollection) GetByEmail(ctx context.Context, email *string) *SGItem {
	return s.Query(ctx, bson.M{"email": *email})
}

// All returns all items in the collection
func (s *SGCollection) All(ctx context.Context) []SGItem {
	var items []SGItem
	cursor, err := s.Mail.Find(ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return nil
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var item SGItem
		err := cursor.Decode(&item)
		if err != nil {
			log.Println(err)
			continue
		}
		items = append(items, item)
	}

	if err := cursor.Err(); err != nil {
		log.Println(err)
	}

	return items
}

// Delete an item in the collection
func (s *SGCollection) Delete(ctx context.Context, item *SGItem) error {
	if _, err := s.Mail.DeleteOne(ctx, bson.DocElem{Name: "_id", Value: item.Subject}); err != nil {
		return err
	}
	return nil
}

// NewSendGrid local store...
func NewSendGrid(db *store.Client) Sendgrid {
	collection := db.DataBase.Collection(sgCollection)
	return &SGCollection{collection}
}
