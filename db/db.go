package db

import (
	"birthdayreminder/config"
	"birthdayreminder/models"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBService struct {
	client             *mongo.Client
	birthdayCollection *mongo.Collection
}

func (d *DBService) InsertBirthday(ctx context.Context, b *models.BirthDay) error {
	result, err := d.birthdayCollection.InsertOne(ctx, b)
	if err != nil {
		return err
	}
	oiD := result.InsertedID.(primitive.ObjectID)
	b.ID = oiD
	return nil
}

func (d *DBService) GetAllBirthDays(ctx context.Context) ([]models.BirthDay, error) {
	var results []models.BirthDay
	cursor, err := d.birthdayCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
func (d *DBService) GetAllBirthDaysByDayMonth(ctx context.Context, day int, month int) ([]models.BirthDay, error) {
	var results []models.BirthDay
	cursor, err := d.birthdayCollection.Find(ctx, bson.M{
		"day":   day,
		"month": month,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func NewDBService(config *config.Config) (*DBService, error) {
	dBURI := fmt.Sprintf("%s:%d", config.DBHost, config.DBPort)
	clientOptions := options.Client().ApplyURI(dBURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}
	birthdayColl := client.Database(config.DBName).Collection("birthdays")
	return &DBService{
		client:             client,
		birthdayCollection: birthdayColl,
	}, nil
}
