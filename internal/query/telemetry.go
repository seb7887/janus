package query

import (
	"fmt"
	"math"
	"strconv"

	"github.com/seb7887/janus/internal/query/builder"
	ts "github.com/seb7887/janus/internal/storage/timescaledb"
	"github.com/seb7887/janus/janusrpc"
	log "github.com/sirupsen/logrus"
)

type QueryServiceTelemetry interface {
	GetTotalSamples(interval string, filters []*janusrpc.Filter) (int, error)
	GetTimeline(req *janusrpc.TimelineQuery) ([]*janusrpc.TimelineResponse, error)
	GetSegments(req *janusrpc.SegmentQuery) ([]*janusrpc.SegmentItem, error)
	GetTimelineSegments(req *janusrpc.SegmentedTimelineQuery) ([]*janusrpc.TimelineResponse, error)
}

type queryServiceTelemetry struct{}

func NewQueryServiceTelemetry() QueryServiceTelemetry {
	return &queryServiceTelemetry{}
}

func (qs *queryServiceTelemetry) GetTotalSamples(interval string, filters []*janusrpc.Filter) (int, error) {
	sql, err := builder.BuildTotalQuery(interval, filters)
	if err != nil {
		return 0, err
	}
	log.Debugf("query: %s", sql)

	return ts.ExecuteTotalQuery(sql)
}

func (qs *queryServiceTelemetry) GetTimeline(req *janusrpc.TimelineQuery) ([]*janusrpc.TimelineResponse, error) {
	sql, err := builder.BuildTimelineQuery(req)
	if err != nil {
		return nil, err
	}
	log.Debugf("query: %s", *sql)

	res, err := ts.ExecuteTMTimelineQuery(*sql, len(req.Dimensions))
	if err != nil {
		return nil, err
	}

	return formatTimelineResponse(req.Dimensions, res)
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
			count, err = roundFloat(v.Count1)
		} else {
			count, err = roundFloat(v.Count2)
		}
		if err != nil {
			return nil, err
		}

		item := &janusrpc.TimelineItem{
			Name:  v.Bucket,
			Count: float32(count),
		}
		items = append(items, item)
	}

	return items, err
}

func (qs *queryServiceTelemetry) GetSegments(req *janusrpc.SegmentQuery) ([]*janusrpc.SegmentItem, error) {
	sql, err := builder.BuildSegmentQuery(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Infof("query: %s", *sql)

	res, err := ts.ExecuteTMSegmentQuery(*sql, len(req.Dimensions))
	if err != nil {
		return nil, err
	}

	return formatSegmentItems(res, len(req.Dimensions))
}

func formatSegmentItems(r []*ts.SegmentQueryResult, numOfDims int) ([]*janusrpc.SegmentItem, error) {
	var segmentItems []*janusrpc.SegmentItem

	if numOfDims < 2 {
		for _, v := range r {
			value, err := roundFloat(v.Count)
			if err != nil {
				return nil, err
			}
			item := &janusrpc.SegmentItem{
				Name:  v.Dim1,
				Value: float32(value),
			}

			segmentItems = append(segmentItems, item)
		}
	} else {
		subsegments, err := getSubSegments(r)
		if err != nil {
			return nil, err
		}
		for k, v := range subsegments {
			item := &janusrpc.SegmentItem{
				Name:     k,
				Segments: v,
			}
			segmentItems = append(segmentItems, item)
		}
	}

	return segmentItems, nil
}

func getSubSegments(r []*ts.SegmentQueryResult) (map[string][]*janusrpc.SegmentItem, error) {
	subsegments := make(map[string][]*janusrpc.SegmentItem)
	var items []*janusrpc.SegmentItem

	for _, v := range r {
		if _, exists := subsegments[v.Dim1]; exists {
			items = subsegments[v.Dim1]
		}

		value, err := roundFloat(v.Count)
		if err != nil {
			return nil, err
		}
		item := &janusrpc.SegmentItem{
			Name:  v.Dim2,
			Value: float32(value),
		}
		items = append(items, item)
		subsegments[v.Dim1] = items
		// Empty slice
		items = []*janusrpc.SegmentItem{}
	}

	return subsegments, nil
}

func (qs *queryServiceTelemetry) GetTimelineSegments(req *janusrpc.SegmentedTimelineQuery) ([]*janusrpc.TimelineResponse, error) {
	sql, err := builder.BuildSegmentedTimelineQuery(req)
	if err != nil {
		return nil, err
	}
	log.Debugf("query: %s", *sql)

	res, err := ts.ExecuteTMSegmentedTimelineQuery(*sql, len(req.GroupBy))
	if err != nil {
		return nil, err
	}

	return formatSegmentedTimelineResponse(res)
}

func formatSegmentedTimelineResponse(r []*ts.SegmentedTimelineResult) ([]*janusrpc.TimelineResponse, error) {
	var res []*janusrpc.TimelineResponse
	itemsBySegmentation := make(map[string][]*janusrpc.TimelineItem)
	var items []*janusrpc.TimelineItem
	var err error

	// First iterate through the array to get all classify items by segmentation
	for _, v := range r {
		count, err := roundFloat(v.Count)
		if err != nil {
			return nil, err
		}

		item := &janusrpc.TimelineItem{
			Name:  v.Bucket,
			Count: float32(count),
		}

		var segmentation string
		if v.Dim2 != "" {
			segmentation = fmt.Sprintf("%s/%s", v.Dim1, v.Dim2)
		} else {
			segmentation = v.Dim1
		}

		if _, exists := itemsBySegmentation[segmentation]; exists {
			items = itemsBySegmentation[segmentation]
		}

		items = append(items, item)
		itemsBySegmentation[segmentation] = items

		// Empty slice
		items = []*janusrpc.TimelineItem{}
	}

	// Iterate through the map to prepare response items
	for k, v := range itemsBySegmentation {
		timelineResponse := &janusrpc.TimelineResponse{
			Dimension: k,
			Items:     v,
		}
		res = append(res, timelineResponse)
	}

	return res, err
}

func roundFloat(str string) (float64, error) {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}
	return math.Round(num*100) / 100, nil
}
