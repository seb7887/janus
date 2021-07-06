package builder

import (
	"fmt"
	"strings"

	"github.com/seb7887/janus/janusrpc"
)

func BuildTimelineQuery(req *janusrpc.TimelineQuery) (*string, error) {
	transformedQuery, err := transformTimelineQuery(req)
	if err != nil {
		return nil, err
	}

	selectClause, err := buildSelectClause(transformedQuery.Granularity, transformedQuery.Interval, transformedQuery.Aggregations)
	if err != nil {
		return nil, err
	}

	whereClause, err := buildWhereClause(transformedQuery.Interval, transformedQuery.Filters)
	if err != nil {
		return nil, err
	}

	groupByClause := buildGroupByClause()
	orderByClause, err := buildOrderByClause(transformedQuery.OrderBy)
	if err != nil {
		return nil, err
	}

	whereClause = fmt.Sprintf("%s %s %s", whereClause, groupByClause, orderByClause)
	searchQuery := buildSearchQuery(selectClause, TELEMETRY_TABLE, whereClause)

	return &searchQuery, nil
}

func transformTimelineQuery(q *janusrpc.TimelineQuery) (*janusrpc.TimelineQuery, error) {
	// First validate time values
	interval := strings.ToUpper(q.Interval)
	if !isValidTimeValue(interval) {
		return nil, fmt.Errorf("Invalid interval %s", interval)
	}

	granularity := strings.ToUpper(q.Granularity)
	if !isValidTimeValue(granularity) {
		return nil, fmt.Errorf("Invalid granularity %s", granularity)
	}

	var aggregations []*janusrpc.Aggregation
	if len(q.Aggregations) == 0 {
		for idx, dimension := range q.Dimensions {
			agg := &janusrpc.Aggregation{
				Type:  "AVG",
				Field: dimension,
				Name:  fmt.Sprintf("count%d", idx+1),
			}
			aggregations = append(aggregations, agg)
		}
	} else {
		aggregations = q.Aggregations
	}

	orderBy := &janusrpc.OrderBy{
		Dimension: "bucket",
		Direction: strings.ToUpper(q.OrderBy.Direction),
	}

	return &janusrpc.TimelineQuery{
		Filters:      q.Filters,
		Dimensions:   q.Dimensions,
		Granularity:  granularity,
		Interval:     interval,
		Aggregations: aggregations,
		OrderBy:      orderBy,
	}, nil
}
