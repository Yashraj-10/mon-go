package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"mon-go/internal/model"
)

const itemsColl = "items"

// ItemStore handles item persistence. Keeps business logic out of the store.
type ItemStore struct {
	coll *mongo.Collection
}

// NewItemStore returns an ItemStore for the given DB.
func NewItemStore(db *DB) *ItemStore {
	return &ItemStore{coll: db.DB.Collection(itemsColl)}
}

// Create inserts a new item and returns it with ID and CreatedAt set.
func (s *ItemStore) Create(ctx context.Context, input model.ItemCreate) (*model.Item, error) {
	doc := model.Item{
		Name:      input.Name,
		CreatedAt: time.Now().UTC(),
	}
	res, err := s.coll.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}
	doc.ID = res.InsertedID.(primitive.ObjectID)
	return &doc, nil
}

// GetByID returns one item by ID, or nil if not found.
func (s *ItemStore) GetByID(ctx context.Context, id primitive.ObjectID) (*model.Item, error) {
	var item model.Item
	err := s.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// List returns items with optional limit (0 = default 100).
func (s *ItemStore) List(ctx context.Context, limit int) ([]model.Item, error) {
	if limit <= 0 {
		limit = 100
	}
	opts := options.Find().SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1})
	cur, err := s.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var items []model.Item
	if err := cur.All(ctx, &items); err != nil {
		return nil, err
	}
	if items == nil {
		items = []model.Item{}
	}
	return items, nil
}

// DeleteByID deletes one item by ID. Returns true if something was deleted.
func (s *ItemStore) DeleteByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
	res, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}
