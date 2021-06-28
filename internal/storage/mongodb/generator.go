package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	generatorColl = "generators"
)

func UpsertGenerator(generator Generator) error {
	client, err := GetMongoClient()
	if err != nil {
		return err
	}

	collection := client.Database(DB).Collection(generatorColl)

	filter := bson.D{primitive.E{Key: "device_id", Value: generator.DeviceId}}
	err = collection.FindOne(context.TODO(), filter).Decode(bson.M{})
	exists := true
	if err != nil {
		if err == mongo.ErrNoDocuments {
			exists = false
		}
		return err
	}

	if !exists {
		// Create new document
		_, err = collection.InsertOne(context.TODO(), generator)
		if err != nil {
			return err
		}
	} else {
		// Update existing document
		updater := bson.M{"$set": bson.M{
			"node_id":          generator.NodeId,
			"temperature":      generator.Temperature,
			"energy_generated": generator.EnergyGenerated,
			"need_manteinance": generator.NeedManteinance,
			"last_manteinance": generator.LastManteinance,
			"enabled":          generator.Enabled,
		}}
		_, err = collection.UpdateOne(context.TODO(), filter, updater)
		if err != nil {
			return err
		}
	}

	return nil
}
