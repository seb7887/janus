package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Meter struct {
	ID             primitive.ObjectID `bson:"_id"`
	DeviceId       string             `bson:"device_id"`
	NodeId         string             `bson:"node_id"`
	Temperature    int64              `bson:"temperature"`
	Consumption    int64              `bson:"consumption"`
	EnergyConsumed int64              `bson:"energy_consumed"`
	LastReport     int64              `bson:"last_report"`
	Connected      bool               `bson:"connected"`
}

type Generator struct {
	ID              primitive.ObjectID `bson:"_id"`
	DeviceId        string             `bson:"device_id"`
	NodeId          string             `bson:"node_id"`
	Temperature     int64              `bson:"temperature"`
	EnergyGenerated int64              `bson:"energy_generated"`
	Enabled         bool               `bson:"enabled"`
	NeedManteinance bool               `bson:"need_manteinance"`
	LastManteinance int64              `bson:"last_manteinance"`
}
