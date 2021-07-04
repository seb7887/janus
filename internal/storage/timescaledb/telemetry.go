package timescaledb

type TimelineQueryResult struct {
	Bucket string
	Count  string
}

func InsertTelemetryEntry(row *Telemetry) error {
	res := DB.Create(&row)
	return res.Error
}

func ExecuteTMTimelineQuery(sql string) ([]*TimelineQueryResult, error) {
	rows, err := DB.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}

	var res []*TimelineQueryResult
	for rows.Next() {
		var (
			bucket string
			count  string
		)
		rows.Scan(&bucket, &count)
		item := &TimelineQueryResult{
			Bucket: bucket,
			Count:  count,
		}
		res = append(res, item)
	}
	return res, nil
}
