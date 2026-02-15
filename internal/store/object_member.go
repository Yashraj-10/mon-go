package store

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"mon-go/internal/model"
)

const objectMemberColl = "object_member"

// ObjectMemberStore handles object_member persistence.
type ObjectMemberStore struct {
	coll *mongo.Collection
}

// NewObjectMemberStore returns an ObjectMemberStore and ensures the unique index on (object_id, member_id).
func NewObjectMemberStore(db *DB) *ObjectMemberStore {
	coll := db.DB.Collection(objectMemberColl)
	ctx := context.Background()
	_, _ = coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "object_id", Value: 1}, {Key: "member_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return &ObjectMemberStore{coll: coll}
}

// CompositeID returns the _id value for an object-member: object_id:member_id.
func CompositeID(objectID, memberID string) string {
	return objectID + ":" + memberID
}

// Create inserts an object-member link with _id = object_id:member_id. Returns duplicate key error if already exists.
func (s *ObjectMemberStore) Create(ctx context.Context, objectID, memberID string) (*model.ObjectMember, error) {
	doc := model.ObjectMember{
		ID:       CompositeID(objectID, memberID),
		ObjectID: objectID,
		MemberID: memberID,
	}
	_, err := s.coll.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// DeleteByID deletes the document with the given _id (object_id:member_id). Returns true if deleted.
func (s *ObjectMemberStore) DeleteByID(ctx context.Context, id string) (bool, error) {
	res, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}

// DeleteByObjectAndMember deletes the document with the given object_id and member_id. Returns true if deleted.
func (s *ObjectMemberStore) DeleteByObjectAndMember(ctx context.Context, objectID, memberID string) (bool, error) {
	res, err := s.coll.DeleteOne(ctx, bson.M{"_id": CompositeID(objectID, memberID)})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}
