package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

// Source represents a news outlet in the database.
// It maps directly to a document in the 'sources' collection.
type Source struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Key              string             `bson:"key" json:"key"`
	Name             string             `bson:"name" json:"name"`
	Description      string             `bson:"description" json:"description"`
	URL              string             `bson:"url" json:"url"`
	LogoURL          string             `bson:"logo_url" json:"logo_url"`
	Languages        []string           `bson:"languages" json:"languages"`
	Topics           []string           `bson:"topics" json:"topics"`
	ReliabilityScore float64            `bson:"reliability_score" json:"reliability_score"`
}
