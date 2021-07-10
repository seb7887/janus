package builder

import (
	"fmt"
	"strings"

	"github.com/seb7887/janus/janusrpc"
)

/* SEGMENTED QUERY
// select device_type, device_id, avg(voltage) from telemetries
where "timestamp" > now() - interval '1 hour'
group by device_type, device_id
*/

func buildSegmentedSelect(dimensions []string, aggregation *janusrpc.Aggregation) (string, error) {
	var clause string
	for _, dim := range dimensions {
		clause = clause + dim + ", "
	}
	aggrExpression, err := buildAggregationsExpression(aggregation)
	if err != nil {
		return clause, err
	}

	return fmt.Sprintf("%s%s", clause, aggrExpression), nil
}

func buildSegmentedGroupBy(dimensions []string) string {
	var clause string
	for idx, dim := range dimensions {
		clause = clause + dim
		if idx < len(dimensions)-1 {
			clause = clause + ", "
		}
	}

	return fmt.Sprintf("GROUP BY %s", clause)
}

func BuildSegmentQuery(req *janusrpc.SegmentQuery) (*string, error) {
	transformedQuery, err := transformSegmentQuery(req)
	if err != nil {
		return nil, err
	}

	selectClause, err := buildSegmentedSelect(transformedQuery.Dimensions, transformedQuery.Aggregation)
	if err != nil {
		return nil, err
	}

	whereClause, err := buildWhereClause(transformedQuery.Interval, transformedQuery.Filters)
	if err != nil {
		return nil, err
	}

	groupByClause := buildSegmentedGroupBy(transformedQuery.Dimensions)

	orderByClause, err := buildOrderByClause(transformedQuery.OrderBy)
	if err != nil {
		return nil, err
	}

	whereClause = fmt.Sprintf("%s %s %s", whereClause, groupByClause, orderByClause)

	searchQuery := buildSearchQuery(selectClause, TELEMETRY_TABLE, whereClause)

	return &searchQuery, nil
}

func transformSegmentQuery(q *janusrpc.SegmentQuery) (*janusrpc.SegmentQuery, error) {
	// First validate time values
	interval := strings.ToUpper(q.Interval)
	if !isValidTimeValue(interval) {
		return nil, fmt.Errorf("Invalid interval %s", interval)
	}

	granularity := strings.ToUpper(q.Granularity)
	if !isValidTimeValue(granularity) {
		return nil, fmt.Errorf("Invalid granularity %s", granularity)
	}

	orderBy := &janusrpc.OrderBy{
		Dimension: q.OrderBy.Dimension,
		Direction: strings.ToUpper(q.OrderBy.Direction),
	}

	aggregation := &janusrpc.Aggregation{
		Name:  q.Aggregation.Name,
		Field: q.Aggregation.Field,
		Type:  strings.ToUpper(q.Aggregation.Type),
	}

	return &janusrpc.SegmentQuery{
		Interval:    interval,
		Granularity: granularity,
		Filters:     q.Filters,
		Dimensions:  q.Dimensions,
		Aggregation: aggregation,
		OrderBy:     orderBy,
	}, nil
}
