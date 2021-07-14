package builder

import (
	"fmt"
	"strings"

	"github.com/seb7887/janus/janusrpc"
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

func buildTimeBucketExpression(granularity string, interval string) string {
	return fmt.Sprintf(`time_bucket_gapfill('%s', timestamp, NOW() - %s, NOW()) as "bucket"`, timeBuckets[granularity], timeRanges[interval])
}

func buildSelectClause(granularity string, interval string, aggregations []*janusrpc.Aggregation) (string, error) {
	timeBucket := buildTimeBucketExpression(granularity, interval)

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

func buildGroupByClause(dimensions []string) string {
	groupByClause := "GROUP BY bucket"
	if len(dimensions) > 0 {
		for _, dim := range dimensions {
			groupByClause = groupByClause + fmt.Sprintf(", %s", dim)
		}
	}

	return groupByClause
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

func BuildTotalQuery(interval string, filters []*janusrpc.Filter) (string, error) {
	whereClause, err := buildWhereClause(strings.ToUpper(interval), filters)
	return fmt.Sprintf(`SELECT COUNT(*) AS "total" FROM %s WHERE %s`, TELEMETRY_TABLE, whereClause), err
}

func isValidTimeValue(str string) bool {
	for _, v := range validTimeValues {
		if v == str {
			return true
		}
	}
	return false
}
