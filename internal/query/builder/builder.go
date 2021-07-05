package builder

import (
	"fmt"
	"strings"

	"github.com/seb7887/janus/janusrpc"
	log "github.com/sirupsen/logrus"
)

const (
	TELEMETRY_TABLE = "telemetries"
)

var (
	validTimeValues = [...]string{"MINUTE", "FIFTEEN_MINUTE", "HOUR", "SIX_HOUR", "DAY", "WEEK", "MONTH", "SEMESTER"}
	timeRanges      = map[string]string{
		"MINUTE":         "INTERVAL '1 minute'",
		"FIFTEEN_MINUTE": "(15* INTERVAL '1 minute')",
		"HOUR":           "INTERVAL '1 hour'",
		"SIX_HOUR":       "(6* INTERVAL '1 hour')",
		"DAY":            "INTERVAL '1 day'",
		"WEEK":           "INTERVAL '1 week'",
		"MONTH":          "INTERVAL '1 month'",
		"SEMESTER":       "(6* INTERVAL '1 month')",
	}
	timeBuckets = map[string]string{
		"MINUTE":         "1 minute",
		"FIFTEEN_MINUTE": "15 minutes",
		"HOUR":           "1 hour",
		"SIX_HOUR":       "6 hours",
		"DAY":            "1 day",
		"WEEK":           "1 week",
		"MONTH":          "1 month",
		"SEMESTER":       "6 months",
	}
)

func buildSearchQuery(selectClause string, dsName string, whereClause string) string {
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s", selectClause, dsName, whereClause)
}

func buildSelectClause(granularity string, aggregations []*janusrpc.Aggregation) (string, error) {
	timeBucket := fmt.Sprintf(`time_bucket('%s', timestamp) as "bucket"`, timeBuckets[granularity])

	var aggregationExpression string
	for idx, v := range aggregations {
		str, err := buildAggregationsExpression(v)
		if err != nil {
			return "", err
		}
		aggregationExpression = aggregationExpression + str
		if idx < len(aggregations)-1 {
			aggregationExpression = aggregationExpression + ", "
		}
	}
	return fmt.Sprintf(`%s, %s`, timeBucket, aggregationExpression), nil
}

func buildAggregationsExpression(aggregation *janusrpc.Aggregation) (string, error) {
	aggregationType := strings.ToUpper(aggregation.Type)
	fieldName := aggregation.Field
	asName := aggregation.Name
	var res string

	switch aggregationType {
	case "COUNT":
		res = fmt.Sprintf("COUNT(%s)", fieldName)
	case "MAX":
		res = fmt.Sprintf("MAX(%s)", fieldName)
	case "MIN":
		res = fmt.Sprintf("MIN(%s)", fieldName)
	case "P25":
		res = fmt.Sprintf("PERCENTILE_COUNT(0.25) within group (order by %s)", fieldName)
	case "P50":
		res = fmt.Sprintf("PERCENTILE_COUNT(0.5) within group (order by %s)", fieldName)
	case "P75":
		res = fmt.Sprintf("PERCENTILE_COUNT(0.75) within group (order by %s)", fieldName)
	case "P90":
		res = fmt.Sprintf("PERCENTILE_COUNT(0.9) within group (order by %s)", fieldName)
	case "P99":
		res = fmt.Sprintf("PERCENTILE_COUNT(0.99) within group (order by %s)", fieldName)
	case "AVG":
		res = fmt.Sprintf("AVG(%s)", fieldName)
	case "SUM":
		res = fmt.Sprintf("SUM(%s)", fieldName)
	default:
		return "", fmt.Errorf("Invalid aggregation type %s", aggregationType)
	}

	return fmt.Sprintf(`%s AS "%s"`, res, asName), nil
}

func buildGroupByClause() string {
	return "GROUP BY bucket"
}

func buildOrderByClause(orderBy *janusrpc.OrderBy) (string, error) {
	if orderBy.Direction != "ASC" && orderBy.Direction != "DESC" {
		return "", fmt.Errorf("Invalid orderBy direction %s", orderBy.Direction)
	}

	return fmt.Sprintf("ORDER BY %s %s", orderBy.Dimension, orderBy.Direction), nil
}

func buildFilter(filter *janusrpc.Filter) (string, error) {
	res := "AND "
	filterType := strings.ToUpper(filter.Type)
	var value string
	if filter.Dimension == "device_id" || filter.Dimension == "node_id" || filter.Dimension == "device_type" {
		value = fmt.Sprintf("'%s'", filter.Value)
	} else {
		value = filter.Value
	}

	switch filterType {
	case "=":
		res = res + fmt.Sprintf("%s = %s", filter.Dimension, value)
	case ">":
		res = res + fmt.Sprintf("%s > %s", filter.Dimension, value)
	case "<":
		res = res + fmt.Sprintf("%s < %s", filter.Dimension, value)
	case ">=":
		res = res + fmt.Sprintf("%s >= %s", filter.Dimension, value)
	case "<=":
		res = res + fmt.Sprintf("%s <= %s", filter.Dimension, value)
	case "!=":
		res = res + fmt.Sprintf("%s NOT %s", filter.Dimension, value)
	case "IN":
		res = res + fmt.Sprintf("%s IN (%s)", filter.Dimension, value)
	case "BETWEEN":
		res = res + fmt.Sprintf("(%s >= %s AND %s <= %s)", filter.Dimension, filter.Lower, filter.Dimension, filter.Upper)
	default:
		return "", fmt.Errorf("Invalid filter type %s", filterType)
	}

	return res, nil
}

func buildWhereClause(interval string, filters []*janusrpc.Filter) (string, error) {
	res := fmt.Sprintf("timestamp > now() - %s", timeRanges[interval])

	if len(filters) > 0 {
		for _, filter := range filters {
			filterClause, err := buildFilter(filter)
			if err != nil {
				return res, err
			}
			res = fmt.Sprintf("%s %s", res, filterClause)
		}
	}

	return res, nil
}

func BuildTotalQuery(req *janusrpc.TimelineQuery) (string, error) {
	whereClause, err := buildWhereClause(strings.ToUpper(req.Interval), req.Filters)
	return fmt.Sprintf(`SELECT COUNT(*) AS "total" FROM %s WHERE %s`, TELEMETRY_TABLE, whereClause), err
}

func BuildTimelineQuery(req *janusrpc.TimelineQuery) (*string, error) {
	transformedQuery, err := transformTimelineQuery(req)
	if err != nil {
		return nil, err
	}

	selectClause, err := buildSelectClause(transformedQuery.Granularity, transformedQuery.Aggregations)
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
	log.Debugf("query: %s", searchQuery)
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

func isValidTimeValue(str string) bool {
	for _, v := range validTimeValues {
		if v == str {
			return true
		}
	}
	return false
}
