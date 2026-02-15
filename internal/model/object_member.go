package model

// ObjectMember is a document in the object_member collection.
// ID is the composite _id: object_id:member_id.
// ObjectID must match group/[0-9]+. MemberID must match group/[0-9]+ or user/[0-9]+.
type ObjectMember struct {
	ID       string `bson:"_id" json:"id"`                           // object_id:member_id
	ObjectID string `bson:"object_id" json:"object_id"`               // group/[0-9]+
	MemberID string `bson:"member_id" json:"member_id"`              // group/[0-9]+ or user/[0-9]+
}

// ObjectMemberCreate is the input for creating an object-member link.
type ObjectMemberCreate struct {
	ObjectID string `json:"object_id"`
	MemberID string `json:"member_id"`
}
