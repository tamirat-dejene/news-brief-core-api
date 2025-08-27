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

func NewSourceRepository(colln *mongo.Collection) *sourceRepository {
	return &sourceRepository{
		collection: colln,
	}
}
func (r *sourceRepository) CreateSource(ctx context.Context, source *entity.Source) error {
	_, err := r.collection.InsertOne(ctx, source)

	if err != nil {
		return err
	}

	return nil
}

// CheckSlugExists checks if a source with the given slug exists.
func (r *sourceRepository) CheckSlugExists(ctx context.Context, slug string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"slug": slug})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckURLExists checks if a source with the given URL exists.
func (r *sourceRepository) CheckURLExists(ctx context.Context, url string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"url": url})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetBySlug retrieves a single source by its unique slug.
func (r *sourceRepository) GetBySlug(ctx context.Context, slug string) (*entity.Source, error) {
	var source entity.Source
	err := r.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&source)
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
