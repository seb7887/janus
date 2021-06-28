package timescaledb

func InsertTelemetryEntry(row *Telemetry) error {
	res := DB.Create(&row)
	return res.Error
}
