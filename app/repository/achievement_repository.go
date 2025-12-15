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
	"github.com/lib/pq" 
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
func (r *AchievementRepository) FindByID(
	ctx context.Context,
	id string,
) (*model.Achievement, error) {

	// 1. Validasi ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid achievement id")
	}

	// 2. Cari berdasarkan _id Mongo
	var ach model.Achievement
	err = r.Coll.FindOne(ctx, bson.M{
		"_id": objID,
	}).Decode(&ach)

	// 3. Tangani kalau data tidak ada
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("achievement not found")
		}
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


func (r *AchievementRepository) SubmitAchievement(
	ctx context.Context,
	mongoAchievementID string,
) error {

	query := `
		UPDATE achievement_references
		SET status = 'submitted',
		    submitted_at = NOW(),
		    updated_at = NOW()
		WHERE mongo_achievement_id = $1
		  AND status = 'draft'
	`

	res, err := r.SqlDB.ExecContext(ctx, query, mongoAchievementID)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("achievement reference not found or already submitted")
	}

	return nil
}

//soft delete
func (r *AchievementRepository) DeleteDraftAchievement(
	ctx context.Context,
	mongoAchievementID string,
) error {

	query := `
		UPDATE achievement_references
		SET status = 'deleted',
		    updated_at = NOW()
		WHERE mongo_achievement_id = $1
		  AND status = 'draft'
	`

	res, err := r.SqlDB.ExecContext(ctx, query, mongoAchievementID)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("achievement not in draft status or not found")
	}

	return nil
}

func (r *AchievementRepository) GetReferencesByStudentIDs(
	ctx context.Context,
	studentIDs []string,
	limit, offset int,
) ([]string, error) {

	query := `
		SELECT mongo_achievement_id
		FROM achievement_references
		WHERE students_id = ANY($1)
		  AND status != 'deleted'
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.SqlDB.QueryContext(
		ctx,
		query,
		pq.Array(studentIDs),
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		rows.Scan(&id)
		ids = append(ids, id)
	}

	return ids, nil
}


func (r *AchievementRepository) FindByIDs(
	ctx context.Context,
	ids []string,
) ([]model.Achievement, error) {

	var objIDs []primitive.ObjectID
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			objIDs = append(objIDs, oid)
		}
	}

	filter := bson.M{
		"_id": bson.M{"$in": objIDs},
	}

	cur, err := r.Coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var achievements []model.Achievement
	for cur.Next(ctx) {
		var ach model.Achievement
		cur.Decode(&ach)
		achievements = append(achievements, ach)
	}

	return achievements, nil
}



func (r *AchievementRepository) ListSubmittedByStudents(
	ctx context.Context,
	studentIDs []string,
) ([]model.Achievement, error) {

	query := `
		SELECT mongo_achievement_id
		FROM achievement_references
		WHERE students_id = ANY($1::uuid[])
		  AND status = 'submitted'
	`

	rows, err := r.SqlDB.QueryContext(ctx, query, pq.Array(studentIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var objIDs []primitive.ObjectID

	for rows.Next() {
		var id string
		rows.Scan(&id)

		oid, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			objIDs = append(objIDs, oid)
		}
	}

	if len(objIDs) == 0 {
		return []model.Achievement{}, nil
	}

	cursor, err := r.Coll.Find(ctx, bson.M{
		"_id": bson.M{"$in": objIDs},
	})
	if err != nil {
		return nil, err
	}

	var achievements []model.Achievement
	err = cursor.All(ctx, &achievements)

	return achievements, err
}


func (r *AchievementRepository) VerifyAchievement(
	ctx context.Context,
	mongoAchievementID string,
	lecturerID string,
) error {

	query := `
		UPDATE achievement_references
		SET status = 'verified',
		    verified_by = $2,
		    verified_at = NOW(),
		    updated_at = NOW()
		WHERE mongo_achievement_id = $1
		  AND status = 'submitted'
	`

	res, err := r.SqlDB.ExecContext(ctx, query, mongoAchievementID, lecturerID)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("achievement not in submitted status or not found")
	}

	return nil
}
