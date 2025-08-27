package mongodb

import (
	"context"
	"errors"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type sourceRepository struct {
	collection *mongo.Collection
}

func NewSourceRepository(db *mongo.Database) *sourceRepository {
	return &sourceRepository{
		collection: db.Collection("sources"),
	}
}

// Exists checks if a source with the given key exists. (Your code was perfect)
func (r *sourceRepository) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"key": key})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetByKey retrieves a single source by its unique key.
func (r *sourceRepository) GetByKey(ctx context.Context, key string) (*entity.Source, error) {
	var source entity.Source
	err := r.collection.FindOne(ctx, bson.M{"key": key}).Decode(&source)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("source not found")
		}
		return nil, err
	}
	return &source, nil
}

// GetAll retrieves all sources from the collection.
func (r *sourceRepository) GetAll(ctx context.Context) ([]entity.Source, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sources []entity.Source
	if err = cursor.All(ctx, &sources); err != nil {
		return nil, err
	}
	return sources, nil
}
