package query

import (
	"github.com/seb7887/janus/internal/query/builder"
	ts "github.com/seb7887/janus/internal/storage/timescaledb"
	"github.com/seb7887/janus/janusrpc"
	log "github.com/sirupsen/logrus"
)

type QueryServiceLog interface {
	GetTotalSamples(interval string, filters []*janusrpc.LogFilter) (int, error)
	GetLogs(req *janusrpc.LogQuery) ([]*janusrpc.LogItem, error)
}

type queryServiceLog struct{}

func NewQueryServiceLog() QueryServiceLog {
	return &queryServiceLog{}
}

func (qs *queryServiceLog) GetTotalSamples(interval string, filters []*janusrpc.LogFilter) (int, error) {
	sql, err := builder.BuildLogTotalQuery(interval, filters)
	if err != nil {
		return 0, err
	}

	log.Debugf("query: %s", *sql)

	return ts.ExecuteTotalQuery(*sql)
}

func (qs *queryServiceLog) GetLogs(req *janusrpc.LogQuery) ([]*janusrpc.LogItem, error) {
	sql, err := builder.BuildLogQuery(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("query: %s", *sql)

	logs, err := ts.ExecuteLogQuery(*sql)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
