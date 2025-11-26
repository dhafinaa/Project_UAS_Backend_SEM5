package repository

import (
	"context"
	"PROJECT_UAS/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository struct {
	Coll *mongo.Collection
}

func (r *AchievementRepository) Create(ctx context.Context, a model.Achievement) (string, error) {
	res, err := r.Coll.InsertOne(ctx, a)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(interface{}).(string), nil
}

func (r *AchievementRepository) FindByID(ctx context.Context, id string) (*model.Achievement, error) {
	var ach model.Achievement

	err := r.Coll.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&ach)

	if err != nil {
		return nil, err
	}

	return &ach, nil
}