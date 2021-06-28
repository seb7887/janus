package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	meterColl = "meters"
)

func UpsertMeter(meter Meter) error {
	client, err := GetMongoClient()
	if err != nil {
		return err
	}

	collection := client.Database(DB).Collection(meterColl)

	filter := bson.D{primitive.E{Key: "device_id", Value: meter.DeviceId}}
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
		_, err = collection.InsertOne(context.TODO(), meter)
		if err != nil {
			return err
		}
	} else {
		// Update existing document
		updater := bson.M{"$set": bson.M{
			"node_id":         meter.NodeId,
			"temperature":     meter.Temperature,
			"consumption":     meter.Consumption,
			"energy_consumed": meter.EnergyConsumed,
			"last_report":     meter.LastReport,
			"connected":       meter.Connected,
		}}
		_, err = collection.UpdateOne(context.TODO(), filter, updater)
		if err != nil {
			return err
		}
	}

	return nil
}
