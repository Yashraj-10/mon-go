package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Item is a sample document for the items collection.
// Add your own fields and collections following this pattern.
type Item struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// ItemCreate is the input for creating an item (no ID, no CreatedAt).
type ItemCreate struct {
	Name string `json:"name"`
}
