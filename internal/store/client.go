package store

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "mail"

// Client wrapper
type Client struct {
	Client   *mongo.Client
	DataBase *mongo.Database
}

// NewMailDB local store...
func NewMailDB(ctx context.Context, dbConn *string) *Client {
	clientOptions := options.Client().ApplyURI(*dbConn)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	return &Client{client, client.Database(dbName)}
}
