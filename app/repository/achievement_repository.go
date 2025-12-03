package repository

import (
	"context"
	"errors"

	"PROJECT_UAS/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository struct {
	Coll *mongo.Collection
}

func NewAchievementRepository(db *mongo.Database) *AchievementRepository {
	return &AchievementRepository{
		Coll: db.Collection("achievements"),
	}
}

func (r *AchievementRepository) Create(ctx context.Context, a model.Achievement) (string, error) {

	res, err := r.Coll.InsertOne(ctx, a)
	if err != nil {
		return "", err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("inserted ID is not ObjectID")
	}

	return oid.Hex(), nil
}

func (r *AchievementRepository) FindByID(ctx context.Context, id string) (*model.Achievement, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var ach model.Achievement

	err = r.Coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&ach)
	if err != nil {
		return nil, err
	}

	return &ach, nil
}
