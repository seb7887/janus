package timescaledb

import (
	"github.com/seb7887/janus/janusrpc"
)

func InsertLogEntry(row *Log) error {
	res := DB.Create(&row)
	return res.Error
}

func ExecuteLogQuery(sql string) ([]*janusrpc.LogItem, error) {
	rows, err := DB.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*janusrpc.LogItem
	for rows.Next() {
		var (
			deviceId  string
			severity  string
			message   string
			timestamp string
		)

		err = rows.Scan(&deviceId, &severity, &message, &timestamp)
		if err != nil {
			return nil, err
		}

		millis, err := getMillis(timestamp)
		if err != nil {
			return nil, err
		}

		item := &janusrpc.LogItem{
			DeviceId:  deviceId,
			Severity:  severity,
			Message:   message,
			Timestamp: millis,
		}

		res = append(res, item)
	}

	return res, nil
}
