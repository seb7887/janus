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

	groupByClause := buildGroupByClause([]string{})
	orderByClause, err := buildOrderByClause(transformedQuery.OrderBy)
	if err != nil {
		return nil, err
	}

	whereClause = fmt.Sprintf("%s %s %s", whereClause, groupByClause, orderByClause)
	searchQuery := buildSearchQuery(selectClause, TELEMETRY_TABLE, whereClause)

	return &searchQuery, nil
}

func sanitizeTimeValue(v string) (string, error) {
	sanitized := strings.ToUpper(v)
	if !isValidTimeValue(sanitized) {
		return sanitized, fmt.Errorf("Invalid time value %s", sanitized)
	}

	return sanitized, nil
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

/* BUCKET RANGE SEGMENTED TIMELINE QUERY
// select time_bucket('15 minutes', timestamp) as "bucket",
case
	when temperature between 0 and 100 then '100'
	when temperature between 100 and 200 then '200'
	else 'nothing'
end as "range"
from telemetries t
where timestamp > now() - interval '1 hour'
group by bucket, temperature
order by bucket asc
*/

func buildGroupBySelectClause(granularity string, interval string, groupBy []string, aggregation *janusrpc.Aggregation) (string, error) {
	timeBucket := buildTimeBucketExpression(granularity, interval)

	var groupByExpression string
	for _, v := range groupBy {
		groupByExpression = groupByExpression + v + ", "
	}

	aggregationExpression, err := buildAggregationsExpression(aggregation)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`%s, %s%s`, timeBucket, groupByExpression, aggregationExpression), nil
}

func BuildSegmentedTimelineQuery(req *janusrpc.SegmentedTimelineQuery) (*string, error) {
	transformedQuery, err := transformSegmentedTimelineQuery(req)
	if err != nil {
		return nil, err
	}

	selectClause, err := buildGroupBySelectClause(transformedQuery.Granularity, transformedQuery.Interval, transformedQuery.GroupBy, transformedQuery.Aggregation)
	if err != nil {
		return nil, err
	}

	whereClause, err := buildWhereClause(transformedQuery.Interval, transformedQuery.Filters)
	if err != nil {
		return nil, err
	}

	groupByClause := buildGroupByClause(transformedQuery.GroupBy)

	orderByClause, err := buildOrderByClause(transformedQuery.OrderBy)
	if err != nil {
		return nil, err
	}

	whereClause = fmt.Sprintf("%s %s %s", whereClause, groupByClause, orderByClause)
	searchQuery := buildSearchQuery(selectClause, TELEMETRY_TABLE, whereClause)

	return &searchQuery, nil
}

func transformSegmentedTimelineQuery(q *janusrpc.SegmentedTimelineQuery) (*janusrpc.SegmentedTimelineQuery, error) {
	interval, err := sanitizeTimeValue(q.Interval)
	if err != nil {
		return nil, err
	}

	granularity, err := sanitizeTimeValue(q.Granularity)
	if err != nil {
		return nil, err
	}

	aggregation := &janusrpc.Aggregation{
		Name:  q.Aggregation.Name,
		Type:  strings.ToUpper(q.Aggregation.Type),
		Field: q.Aggregation.Field,
	}

	orderBy := &janusrpc.OrderBy{
		Dimension: "bucket",
		Direction: strings.ToUpper(q.OrderBy.Direction),
	}

	return &janusrpc.SegmentedTimelineQuery{
		Filters:      q.Filters,
		Granularity:  granularity,
		Interval:     interval,
		Aggregation:  aggregation,
		GroupBy:      q.GroupBy,
		BucketRanges: q.BucketRanges,
		OrderBy:      orderBy,
	}, nil
}
