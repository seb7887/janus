package query

import (
	"math"
	"strconv"

	"github.com/seb7887/janus/internal/query/builder"
	ts "github.com/seb7887/janus/internal/storage/timescaledb"
	"github.com/seb7887/janus/janusrpc"
)

type QueryServiceTelemetry interface {
	GetTotalSamples(req *janusrpc.TimelineQuery) (int, error)
	GetTimeline(req *janusrpc.TimelineQuery) ([]*janusrpc.TimelineResponse, error)
}

type queryServiceTelemetry struct{}

func NewQueryServiceTelemetry() QueryServiceTelemetry {
	return &queryServiceTelemetry{}
}

func (qs *queryServiceTelemetry) GetTotalSamples(req *janusrpc.TimelineQuery) (int, error) {
	sql, err := builder.BuildTotalQuery(req)
	if err != nil {
		return 0, err
	}

	return ts.ExecuteTMTotalQuery(sql)
}

func (qs *queryServiceTelemetry) GetTimeline(req *janusrpc.TimelineQuery) ([]*janusrpc.TimelineResponse, error) {
	sql, err := builder.BuildTimelineQuery(req)
	if err != nil {
		return nil, err
	}

	res, err := ts.ExecuteTMTimelineQuery(*sql, len(req.Dimensions))
	if err != nil {
		return nil, err
	}

	formatedResponse, err := formatTimelineResponse(req.Dimensions, res)

	return formatedResponse, err
}

func formatTimelineResponse(dims []string, r []*ts.TimelineQueryResult) ([]*janusrpc.TimelineResponse, error) {
	var res []*janusrpc.TimelineResponse

	for idx, dim := range dims {
		dimItems, err := formatTimelineItems(r, idx+1)
		if err != nil {
			return nil, err
		}

		dimRes := &janusrpc.TimelineResponse{
			Dimension: dim,
			Items:     dimItems,
		}

		res = append(res, dimRes)
	}

	return res, nil
}

func formatTimelineItems(r []*ts.TimelineQueryResult, idx int) ([]*janusrpc.TimelineItem, error) {
	var items []*janusrpc.TimelineItem
	var err error
	for _, v := range r {
		var count float64
		if idx == 1 {
			count, err = strconv.ParseFloat(v.Count1, 64)
		} else {
			count, err = strconv.ParseFloat(v.Count2, 64)
		}
		if err != nil {
			return nil, err
		}
		rounded := math.Round(count*100) / 100

		item := &janusrpc.TimelineItem{
			Name:  v.Bucket,
			Count: float32(rounded),
		}
		items = append(items, item)
	}

	return items, err
}
