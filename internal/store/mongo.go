package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB wraps the MongoDB client and database for dependency injection and testing.
type DB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// NewMongoDB connects to MongoDB and returns a DB instance.
func NewMongoDB(ctx context.Context, uri, dbName string) (*DB, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, err
	}

	return &DB{
		Client: client,
		DB:     client.Database(dbName),
	}, nil
}

// Close disconnects the MongoDB client.
func (db *DB) Close(ctx context.Context) error {
	return db.Client.Disconnect(ctx)
}
