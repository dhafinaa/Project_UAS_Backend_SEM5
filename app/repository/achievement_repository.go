package repository

import (
	"context"
	"errors"
	"time"

	"PROJECT_UAS/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository struct {
	Coll      *mongo.Collection
	RefColl   *mongo.Collection
}

func NewAchievementRepository(db *mongo.Database) *AchievementRepository {
	return &AchievementRepository{
		Coll:    db.Collection("achievements"),
		RefColl: db.Collection("achievement_references"),
	}
}

//
// -----------------------------------------------------------
// CREATE ACHIEVEMENT
// -----------------------------------------------------------
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

//
// -----------------------------------------------------------
// FIND BY ID
// -----------------------------------------------------------
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

//
// -----------------------------------------------------------
// FIND ACHIEVEMENTS BY STUDENT ID
// -----------------------------------------------------------
func (r *AchievementRepository) FindByStudentID(ctx context.Context, studentID string) ([]model.Achievement, error) {

	cur, err := r.Coll.Find(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.Achievement

	for cur.Next(ctx) {
		var a model.Achievement
		cur.Decode(&a)
		list = append(list, a)
	}

	return list, nil
}

//
// -----------------------------------------------------------
// DELETE ACHIEVEMENT BY ID
// -----------------------------------------------------------
func (r *AchievementRepository) DeleteByID(ctx context.Context, id string) error {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.Coll.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

//
// -----------------------------------------------------------
// CREATE REFERENCE (SQL replacement, stored in Mongo since your model is Mongo)
// -----------------------------------------------------------
func (r *AchievementRepository) CreateReference(ctx context.Context, ref model.Achievement_reference) error {

	_, err := r.RefColl.InsertOne(ctx, ref)
	return err
}

// LIST ACHIEVEMENTS BY STUDENT ID
func (r *AchievementRepository) ListByStudent(ctx context.Context, studentID string) ([]model.Achievement, error) {

	cur, err := r.Coll.Find(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.Achievement

	for cur.Next(ctx) {
		var a model.Achievement
		cur.Decode(&a)
		list = append(list, a)
	}

	return list, nil
}

// UPDATE STATUS IN REFERENCE COLLECTION
func (r *AchievementRepository) UpdateStatusByID(ctx context.Context, achID string, status string) error {

	_, err := r.RefColl.UpdateOne(
		ctx,
		bson.M{"mongo_achievement_id": achID},
		bson.M{
			"$set": bson.M{
				"status":      status,
				"updated_at":  time.Now(),
				"verified_at": time.Now(),
			},
		},
	)

	return err
}
