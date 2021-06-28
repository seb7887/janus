package timescaledb

func InsertLogEntry(row *Log) error {
	res := DB.Create(&row)
	return res.Error
}
