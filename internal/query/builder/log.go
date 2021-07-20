package builder

import (
	"fmt"
	"strings"

	"github.com/seb7887/janus/janusrpc"
)

var (
	TEMPLATE = `SELECT device_id, severity, message, "timestamp" FROM "%s" WHERE %s ORDER BY "timestamp" DESC %s`
)

func buildLogFilterClause(interval string, filters []*janusrpc.LogFilter) string {
	intervalExpression := buildTimeInterval(interval)

	var filterExpression string
	for _, filter := range filters {
		filterExpression = filterExpression + fmt.Sprintf("%s = '%s'", filter.Field, filter.Value) + " AND "
	}

	return fmt.Sprintf("%s%s", filterExpression, intervalExpression)
}

func buildPaginationClause(limit int64, offset int64) string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

func BuildLogTotalQuery(interval string, filters []*janusrpc.LogFilter) (*string, error) {
	whereClause := buildLogFilterClause(strings.ToUpper(interval), filters)
	searchQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", LOG_TABLE, whereClause)

	return &searchQuery, nil
}

func BuildLogQuery(q *janusrpc.LogQuery) (*string, error) {
	if q.Limit == 0 || q.Offset < 0 {
		return nil, fmt.Errorf("Invalid limit or offset values")
	}

	whereClause := buildLogFilterClause(strings.ToUpper(q.Interval), q.Filters)
	paginationClause := buildPaginationClause(q.Limit, q.Offset)

	searchQuery := fmt.Sprintf(TEMPLATE, LOG_TABLE, whereClause, paginationClause)

	return &searchQuery, nil
}
