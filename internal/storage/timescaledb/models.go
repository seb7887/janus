package timescaledb

import (
	"time"
)

type Telemetry struct {
	DeviceId    string
	DeviceType  string
	Power       int64
	Voltage     int64
	Current     int64
	Temperature int64
	Timestamp   time.Time
}

type Log struct {
	DeviceId  string
	Subject   string
	Message   string
	Timestamp time.Time
}
