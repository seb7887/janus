package builder

import (
	"fmt"
	"strings"

	"github.com/seb7887/janus/janusrpc"
)

func BuildTimelineQuery(req *janusrpc.TimelineQuery) (*string, error) {
	q, err := transformTimelineQuery(req)
	if err != nil {
		return nil, err
	}

	selectClause, err := buildSelectClause(q.Granularity, q.Interval, q.Aggregations)
	if err != nil {
		return nil, err
	}

	whereClause, err := buildWhereClause(q.Interval, q.Filters)
	if err != nil {
		return nil, err
	}

	groupByClause := buildGroupByClause([]string{})
	orderByClause, err := buildOrderByClause(q.OrderBy)
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
	interval, err := sanitizeTimeValue(q.Interval)
	if err != nil {
		return nil, err
	}

	granularity, err := sanitizeTimeValue(q.Granularity)
	if err != nil {
		return nil, err
	}

	if len(q.Dimensions) > 2 {
		return nil, fmt.Errorf("Invalid number of dimensions")
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

func buildBucketRangeExpression(buckets *janusrpc.SegmentBucket, name string) string {
	var conditions string

	for _, b := range buckets.BucketRanges {
		conditions = conditions + fmt.Sprintf("WHEN %s BETWEEN %s AND %s THEN '%s' ", buckets.Dimension, b.Lower, b.Upper, b.Name)
	}
	conditions = conditions + "ELSE 'unknown'"

	return fmt.Sprintf(`CASE %s END AS "%s"`, conditions, name)
}

func buildGroupBySelectClause(granularity string, interval string, groupBy []string, aggregation *janusrpc.Aggregation, buckets *janusrpc.SegmentBucket) (string, error) {
	timeBucket := buildTimeBucketExpression(granularity, interval)

	var groupByExpression string
	if len(buckets.BucketRanges) > 0 {
		groupByExpression = buildBucketRangeExpression(buckets, "segment") + ", "
	}

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
	q, err := transformSegmentedTimelineQuery(req)
	if err != nil {
		return nil, err
	}

	selectClause, err := buildGroupBySelectClause(q.Granularity, q.Interval, q.GroupBy, q.Aggregation, q.SegmentBucket)
	if err != nil {
		return nil, err
	}

	whereClause, err := buildWhereClause(q.Interval, q.Filters)
	if err != nil {
		return nil, err
	}

	groupBy := q.GroupBy
	if len(q.SegmentBucket.BucketRanges) > 0 {
		groupBy = append(groupBy, q.SegmentBucket.Dimension)
	}
	groupByClause := buildGroupByClause(groupBy)

	orderByClause, err := buildOrderByClause(q.OrderBy)
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

	if len(q.GroupBy) > 2 || (len(q.SegmentBucket.BucketRanges) > 0 && len(q.GroupBy) == 2) {
		return nil, fmt.Errorf("Invalid number of dimensions")
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
		Filters:       q.Filters,
		Granularity:   granularity,
		Interval:      interval,
		Aggregation:   aggregation,
		GroupBy:       q.GroupBy,
		SegmentBucket: q.SegmentBucket,
		OrderBy:       orderBy,
	}, nil
}
