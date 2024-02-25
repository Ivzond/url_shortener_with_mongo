package app

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UrlDAO struct {
	c *mongo.Collection
}

func NewUrlDAO(ctx context.Context, client *mongo.Client) (*UrlDAO, error) {
	dao := &UrlDAO{
		c: client.Database("core").Collection("shortUrls"),
	}
	if err := dao.createIndices(ctx); err != nil {
		return nil, err
	}
	return dao, nil
}

func (dao *UrlDAO) createIndices(ctx context.Context) error {
	_, err := dao.c.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{"expireAt", 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	return err
}

func (dao *UrlDAO) Insert(ctx context.Context, shortUrl *ShortURL) error {
	_, err := dao.c.InsertOne(ctx, shortUrl)
	return err
}

func (dao *UrlDAO) FindByID(ctx context.Context, id string) (*ShortURL, error) {
	filter := bson.D{{"_id", id}}
	var shortURL ShortURL
	err := dao.c.FindOne(ctx, filter).Decode(&shortURL)
	switch {
	case err == nil:
		return &shortURL, nil
	case errors.Is(err, mongo.ErrNoDocuments):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (dao *UrlDAO) Update(ctx context.Context, url *ShortURL) error {
	filter := bson.D{{"_id", url.ID}}
	result, err := dao.c.ReplaceOne(ctx, filter, url)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (dao *UrlDAO) DeleteByID(ctx context.Context, id string) error {
	filter := bson.D{{"_id", id}}
	result, err := dao.c.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}
