package timescaledb

import (
	"fmt"
	"time"
)

type TimelineQueryResult struct {
	Bucket string
	Count1 string
	Count2 string
}

func InsertTelemetryEntry(row *Telemetry) error {
	res := DB.Create(&row)
	return res.Error
}

func ExecuteTMTotalQuery(sql string) (int, error) {
	var total int
	rows, err := DB.Raw(sql).Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			return 0, err
		}
	}
	return total, nil
}

func ExecuteTMTimelineQuery(sql string, numOfDims int) ([]*TimelineQueryResult, error) {
	rows, err := DB.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*TimelineQueryResult
	for rows.Next() {
		var (
			bucket string
			count1 []uint8
			count2 []uint8
		)

		if numOfDims < 2 {
			err = rows.Scan(&bucket, &count1)
		} else {
			err = rows.Scan(&bucket, &count1, &count2)
		}
		if err != nil {
			return nil, err
		}

		bucketMillis, err := getMillis(bucket)
		if err != nil {
			return nil, err
		}

		item := &TimelineQueryResult{
			Bucket: fmt.Sprintf("%d", bucketMillis),
			Count1: convertUint8ToStr(count1),
			Count2: convertUint8ToStr(count2),
		}
		res = append(res, item)
	}

	return res, nil
}

func convertUint8ToStr(u []uint8) string {
	return fmt.Sprintf("%s", u)
}

func getMillis(str string) (int64, error) {
	t, err := time.Parse(time.RFC3339, str)
	millis := t.UTC().UnixNano() / 1000000

	return millis, err
}
