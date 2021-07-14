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

type SegmentQueryResult struct {
	Dim1  string
	Dim2  string
	Count string
}

type SegmentedTimelineResult struct {
	Bucket string
	Dim1   string
	Dim2   string
	Count  string
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

func ExecuteTMSegmentQuery(sql string, numOfDims int) ([]*SegmentQueryResult, error) {
	rows, err := DB.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*SegmentQueryResult
	for rows.Next() {
		var (
			dim1  string
			dim2  string
			count []uint8
		)

		if numOfDims < 2 {
			err = rows.Scan(&dim1, &count)
		} else {
			err = rows.Scan(&dim1, &dim2, &count)
		}
		if err != nil {
			return nil, err
		}

		item := &SegmentQueryResult{
			Dim1:  dim1,
			Dim2:  dim2,
			Count: convertUint8ToStr(count),
		}
		res = append(res, item)
	}

	return res, nil
}

func ExecuteTMSegmentedTimelineQuery(sql string, numOfDims int) ([]*SegmentedTimelineResult, error) {
	rows, err := DB.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*SegmentedTimelineResult
	for rows.Next() {
		var (
			bucket string
			dim1   string
			dim2   string
			count  []uint8
		)

		if numOfDims < 2 {
			err = rows.Scan(&bucket, &dim1, &count)
		} else {
			err = rows.Scan(&bucket, &dim1, &dim2, &count)
		}
		if err != nil {
			return nil, err
		}

		bucketMillis, err := getMillis(bucket)
		if err != nil {
			return nil, err
		}

		item := &SegmentedTimelineResult{
			Bucket: fmt.Sprintf("%d", bucketMillis),
			Dim1:   dim1,
			Dim2:   dim2,
			Count:  convertUint8ToStr(count),
		}
		res = append(res, item)
	}

	return res, nil
}

func convertUint8ToStr(u []uint8) string {
	if len(u) == 0 {
		return "0"
	}
	return fmt.Sprintf("%s", u)
}

func getMillis(str string) (int64, error) {
	t, err := time.Parse(time.RFC3339, str)
	millis := t.UTC().UnixNano() / 1000000

	return millis, err
}
