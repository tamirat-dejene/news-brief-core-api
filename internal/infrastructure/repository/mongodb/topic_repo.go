package mongodb

import (
	"context"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type topicRepository struct {
	collection *mongo.Collection
}

// NewTopicRepository creates a new MongoDB topic repository.
func NewTopicRepository(db *mongo.Database) *topicRepository {
	return &topicRepository{
		collection: db.Collection("topics"),
	}
}

// GetAll retrieves all topics from the 'topics' collection. // CORRECTED COMMENT
func (r *topicRepository) GetAll(ctx context.Context) ([]entity.Topic, error) {
	var topics []entity.Topic

	// Sort by 'sort_order' as defined in the schema file
	opts := options.Find().SetSort(bson.D{{Key: "sort_order", Value: 1}})
	cursor, err := r.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &topics); err != nil {
		return nil, err
	}

	return topics, nil
}
