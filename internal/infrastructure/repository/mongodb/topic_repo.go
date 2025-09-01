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
func NewTopicRepository(colln *mongo.Collection) *topicRepository {
	return &topicRepository{
		collection: colln,
	}
}

func (r *topicRepository) CreateTopic(ctx context.Context, topic *entity.Topic) error {
	_, err := r.collection.InsertOne(ctx, topic)
	if err != nil {
		return err
	}
	return nil
}

func (r *topicRepository) CheckSlugExists(ctx context.Context, slug string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"slug": slug})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAll retrieves all topics from the 'topics' collection.
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
