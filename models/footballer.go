package models

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Footballer struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	FootballClub string `json:"football_club"`
}

func GetAllFootballers() ([]*Footballer, error) {
	filter := bson.D{{}}
	return filterFootballers(filter)
}

func GetFootballerByID(id string) (*Footballer, error) {
	filter := bson.M{"account_id": id}
	var object Footballer

	if err := footballerCollection.FindOne(context.Background(), filter).Decode(&object); err != nil {
		return nil, err
	}
	return &object, nil
}

func CreateFootballer(footballer *Footballer) (*Footballer, error) {
	footballer.ID = uuid()
	_, err := footballerCollection.InsertOne(ctx, footballer)
	return footballer, err
}

func UpdateFootballer(id string, updatedFootballer *Footballer) (*Footballer, error) {
	filter := bson.M{"footballer_id": id}
	update := bson.M{"$set": bson.M{
		"first_name":      updatedFootballer.FirstName,
		"last_name":       updatedFootballer.LastName,
		"footballer_club": updatedFootballer.FootballClub,
	}}

	result, err := footballerCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		panic(err)
	}
	fmt.Println("modified count: ", result.ModifiedCount)

	footballer, err := GetFootballerByID(id)
	if err != nil {
		panic(err)
	}
	return footballer, err
}

func DeleteFootballer(id string) (bool, error) {
	footballer, err := GetFootballerByID(id)
	if err != nil {
		panic(err)
	}

	footballerCollection.DeleteOne(ctx, bson.M{"id": footballer.ID})

	if err != nil {
		log.Fatal("DeleteFootballer ERROR:", err)
		return false, err
	}

	return true, err
}

func filterFootballers(filter interface{}) ([]*Footballer, error) {
	var footballers []*Footballer

	cur, err := footballerCollection.Find(ctx, filter)
	if err != nil {
		return footballers, err
	}

	for cur.Next(ctx) {
		var f Footballer
		err := cur.Decode(&f)
		if err != nil {
			return footballers, err
		}

		footballers = append(footballers, &f)
	}

	if err := cur.Err(); err != nil {
		return footballers, err
	}

	cur.Close(ctx)

	if len(footballers) == 0 {
		return footballers, mongo.ErrNoDocuments
	}

	return footballers, nil
}

func uuid() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)

	if err != nil {
		return ""
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
