package eq

import (
	"strconv"
	"strings"
)

var (
	// Question 占位符作为 (?)
	Question = Dialect{Placeholder: "?"}

	// Dollar 占位符为 ($1, $2, $3)
	Dollar = Dialect{Placeholder: "$"}

	// Colon 占位符为 (:1, :2, :3）
	Colon = Dialect{Placeholder: ":"}

	// AtP 占位符为 (@p1, @p2, @p3)
	AtP = Dialect{Placeholder: "@"}

	DefaultPlaceholder = Question
)

type Placeholder interface {
	ReplacePlaceholders(sql string) string
}

type Dialect struct {
	Placeholder string
}

func (d Dialect) Where(a ...Builder) *List {
	return Where(a...).Placeholder(d)
}

func (d Dialect) And(a ...Builder) *List {
	return And(a...).Placeholder(d)
}

func (d Dialect) Or(a ...Builder) *List {
	return Or(a...).Placeholder(d)
}

func (d Dialect) OrLine(a ...Builder) *List {
	return OrLine(a...).Placeholder(d)
}

func (d Dialect) ReplacePlaceholders(sql string) string {
	if d.Placeholder == "?" {
		return sql
	}
	return ReplacePositionalPlaceholders(sql, d.Placeholder)
}

func ReplacePositionalPlaceholders(sql, prefix string) string {
	sqls := strings.Split(sql, "?")
	if len(sqls) == 1 {
		return sql
	}

	var sb strings.Builder
	for i, x := range sqls[:len(sqls)-1] {
		sb.WriteString(x + prefix)
		sb.WriteString(strconv.Itoa(i + 1))
	}

	sb.WriteString(sqls[len(sqls)-1])
	return sb.String()
}
