// sqb package is a simple SQL query builder that allows you to build SQL queries in a more readable way.
package sqb

import (
	"fmt"
	"strings"
	"unicode"
)

type SQLQueryBuilder struct {
	SelectedColumns  []string
	WhereClauses     []string
	LeftJoins        []string
	InnerJoins       []string
	FromTable        string
	OrderByColumn    string
	OrderByDirection string

	// Named parameters in the query. For example, $userId, $a_b_c
	Parameters map[string]any
}

func NewSQLQueryBuilder() *SQLQueryBuilder {
	return &SQLQueryBuilder{
		Parameters: map[string]any{},
	}
}

// Build returns the SQL query and the parameters to be used in the query.
func (sqb *SQLQueryBuilder) Build() (string, []interface{}) {
	query := strings.Builder{}
	if len(sqb.SelectedColumns) == 0 {
		query.WriteString("SELECT *")
	} else {
		query.WriteString("SELECT ")
		query.WriteString(strings.Join(sqb.SelectedColumns, ", "))
	}

	query.WriteString(" FROM ")
	query.WriteString(sqb.FromTable)

	if len(sqb.LeftJoins) > 0 {
		query.WriteString(" LEFT JOIN ")
		query.WriteString(strings.Join(sqb.LeftJoins, " LEFT JOIN "))
	}

	if len(sqb.InnerJoins) > 0 {
		query.WriteString(" INNER JOIN ")
		query.WriteString(strings.Join(sqb.InnerJoins, " INNER JOIN "))
	}

	if len(sqb.WhereClauses) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(sqb.WhereClauses, " AND "))
	}

	if sqb.OrderByColumn != "" {
		query.WriteString(" ORDER BY ")
		query.WriteString(sqb.OrderByColumn)
		query.WriteString(" ")
		query.WriteString(sqb.OrderByDirection)
	}

	// Parsing named parameters in the query and replacing it by their numeric equivalent.
	// For example:
	// SELECT * FROM orders WHERE customer_id = $customerId
	// Output: SELECT * FROM orders WHERE customer_id = $1
	raw := query.String()
	params := parseParameters(raw)
	args := make([]interface{}, 0, len(params))
	for _, key := range params {
		if value, ok := sqb.Parameters[key]; ok {
			args = append(args, value)
		}
	}
	replaced := replaceParamsToSQLVars(raw, params)

	return replaced, args
}

func isIdentifier(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_'
}

// Parses variables in the query, like $abc and returns an array of the params without '$' symbol
func parseParameters(query string) []string {
	exists := map[string]struct{}{}
	args := []string{}
	for i := 0; i < len(query); i++ {
		if query[i] == '$' {
			i++
			start := i
			for i < len(query) && isIdentifier(rune(query[i])) {
				i++
			}
			arg := query[start:i]
			if _, ok := exists[arg]; !ok && arg != "" {
				args = append(args, arg)
				exists[arg] = struct{}{}
			}
		}
	}

	return args
}

// This function replaces all named parameters ($abc, $test) to their numeric equivalent ($1, $2, $3, ...$n)
func replaceParamsToSQLVars(query string, args []string) string {
	for i, arg := range args {
		arg := strings.Trim(arg, " ")
		if arg != "" {
			query = strings.ReplaceAll(query, "$"+arg, fmt.Sprintf("$%d", i+1))
		}
	}
	return query
}

func (sqb *SQLQueryBuilder) Select(columns ...string) *SQLQueryBuilder {
	sqb.SelectedColumns = columns
	return sqb
}

func (sqb *SQLQueryBuilder) From(table string) *SQLQueryBuilder {
	sqb.FromTable = table
	return sqb
}

func (sqb *SQLQueryBuilder) Where(clause string) *SQLQueryBuilder {
	sqb.WhereClauses = []string{clause}
	return sqb
}

func (sqb *SQLQueryBuilder) AndWhere(clause string) *SQLQueryBuilder {
	sqb.WhereClauses = append(sqb.WhereClauses, clause)
	return sqb
}

func (sqb *SQLQueryBuilder) LeftJoin(table, condition string) *SQLQueryBuilder {
	sqb.LeftJoins = append(sqb.LeftJoins, table+" ON "+condition)
	return sqb
}

func (sqb *SQLQueryBuilder) InnerJoin(table, condition string) *SQLQueryBuilder {
	sqb.InnerJoins = append(sqb.InnerJoins, table+" ON "+condition)
	return sqb
}

func (sqb *SQLQueryBuilder) OrderBy(column, direction string) *SQLQueryBuilder {
	sqb.OrderByColumn = column
	sqb.OrderByDirection = direction
	return sqb
}

func (sqb *SQLQueryBuilder) SetParameter(key string, value any) *SQLQueryBuilder {
	sqb.Parameters[key] = value
	return sqb
}
