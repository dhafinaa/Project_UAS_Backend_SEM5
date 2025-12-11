package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"PROJECT_UAS/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository struct {
	Coll  *mongo.Collection // MongoDB achievements
	SqlDB *sql.DB           // PostgreSQL achievement_references
}

func NewAchievementRepository(mongo *mongo.Database, sql *sql.DB) *AchievementRepository {
	return &AchievementRepository{
		Coll:  mongo.Collection("achievements"),
		SqlDB: sql,
	}
}

//
// CREATE ACHIEVEMENT (MongoDB)
//
func (r *AchievementRepository) Create(ctx context.Context, a model.Achievement) (string, error) {
	res, err := r.Coll.InsertOne(ctx, a)
	if err != nil {
		return "", err
	}

	oid := res.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}

//
// FIND BY ID
//
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
// FIND BY STUDENT ID
//
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
// ALIAS FOR LECTURER SERVICE
//
func (r *AchievementRepository) ListByStudent(ctx context.Context, studentID string) ([]model.Achievement, error) {
	return r.FindByStudentID(ctx, studentID)
}

//
// DELETE (Mongo)
//
func (r *AchievementRepository) DeleteByID(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.Coll.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

//
// CREATE REFERENCE (PostgreSQL)
//
func (r *AchievementRepository) CreateReference(ctx context.Context, studentID string, mongoAchievementID string) error {
    if r.SqlDB == nil {
        return errors.New("sql db not configured")
    }

    query := `
        INSERT INTO achievement_references
        (students_id, mongo_achievement_id, status, submitted_at, created_at, updated_at)
        VALUES ($1, $2, 'draft', NOW(), NOW(), NOW())
    `

    // DEBUG
    fmt.Println("DEBUG INSERT studentID =", studentID)
    fmt.Println("DEBUG INSERT mongoID =", mongoAchievementID)

    _, err := r.SqlDB.ExecContext(ctx, query, studentID, mongoAchievementID)
    
    if err != nil {
        fmt.Println("SQL ERROR:", err)  // <---- INI YANG PALING PENTING
    }

    return err
}


//
// UPDATE STATUS (PostgreSQL)
//
func (r *AchievementRepository) UpdateStatusByID(ctx context.Context, achID string, status string) error {
	if r.SqlDB == nil {
		return errors.New("sql db not configured")
	}

	query := `
		UPDATE achievement_references
		SET status=$1, updated_at=NOW(), verified_at=NOW()
		WHERE mongo_achievement_id=$2
	`

	_, err := r.SqlDB.ExecContext(ctx, query, status, achID)
	return err
}


// -----------------------------------------------------------
// UPDATE ACHIEVEMENT (MongoDB)
// -----------------------------------------------------------
func (r *AchievementRepository) UpdateAchievement(ctx context.Context, id string, update bson.M) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.Coll.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": update},
	)
	return err
}

// -----------------------------------------------------------
// ADD ATTACHMENT
// -----------------------------------------------------------
func (r *AchievementRepository) AddAttachment(ctx context.Context, id string, attachment model.Attachment) error {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    _, err = r.Coll.UpdateOne(
        ctx,
        bson.M{"_id": objID},
        bson.M{
            "$push": bson.M{"attachments": attachment},
            "$set":  bson.M{"updated_at": time.Now()},
        },
    )
    return err
}
