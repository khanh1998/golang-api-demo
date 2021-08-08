package model

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Ingredient struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name string             `bson:"name,omitempty" json:"name"`
	Unit string             `bson:"unit,omitempty" json:"unit"`
	V    int                `bson:"__v,omitempty" json:"__v"`
}

type Ingredients struct {
	collection *mongo.Collection
}

func NewIngredients(collectionName string, database *mongo.Database) *Ingredients {
	ingredient := database.Collection(collectionName)
	return &Ingredients{
		collection: ingredient,
	}
}

func (i *Ingredients) GetAll() ([]Ingredient, error) {
	var ingredients []Ingredient
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := i.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Panic(err)
	}
	if err = cursor.All(ctx, &ingredients); err != nil {
		log.Panic(err)
	}
	return ingredients, err
}

func (i *Ingredients) GetById(id string) (Ingredient, error) {
	var ingredient Ingredient
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Panic(err)
	}
	filter := bson.M{"_id": objectId}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := i.collection.FindOne(ctx, filter).Decode(&ingredient); err != nil {
		log.Panic(err)
	}
	return ingredient, err
}

func (i *Ingredients) find(filter interface{}) []Ingredient {
	return []Ingredient{}
}

func (i *Ingredients) AddOne(data Ingredient) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	response, err := i.collection.InsertOne(ctx, data)
	if err != nil {
		log.Panic(err)
	}
	if objectId, ok := response.InsertedID.(primitive.ObjectID); ok {
		return objectId.Hex(), nil
	} else {
		return "", errors.New("cannot convert object id")
	}
}

func (i *Ingredients) UpdateOne(id string, data Ingredient) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Panic(err)
	}
	ingredient := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "name", Value: data.Name},
			primitive.E{Key: "unit", Value: data.Unit},
		}},
	}
	response, err := i.collection.UpdateByID(ctx, objectId, ingredient)
	success := response.ModifiedCount == 1
	return success, nil
}

func (i *Ingredients) DeleteOne(id string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Panic(err)
	}
	filter := bson.M{"_id": objectId}
	response, err := i.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Panic(err)
	}
	return response.DeletedCount == 1
}
