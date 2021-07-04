package query

import (
	"github.com/seb7887/janus/internal/query/builder"
	ts "github.com/seb7887/janus/internal/storage/timescaledb"
	"github.com/seb7887/janus/janusrpc"
)

type QueryServiceTelemetry interface {
	GetTimeline(req *janusrpc.TimelineQuery) error
}

type queryServiceTelemetry struct{}

func NewQueryServiceTelemetry() QueryServiceTelemetry {
	return &queryServiceTelemetry{}
}

func (qs *queryServiceTelemetry) GetTimeline(req *janusrpc.TimelineQuery) error {
	sql, err := builder.BuildTimelineQuery(req)
	_, err = ts.ExecuteTMTimelineQuery(*sql)

	return err
}
