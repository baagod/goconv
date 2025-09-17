package eq

import (
	"bytes"
	"fmt"
	"strings"
)

type Placeholder interface {
	ReplacePlaceholders(sql string) string
}

type Dialect struct {
	Formatter Placeholder
}

func (d Dialect) Where(a ...Builder) *List {
	return Where(a...).Placeholder(d.Formatter)
}

func (d Dialect) And(a ...Builder) *List {
	return And(a...).Placeholder(d.Formatter)
}

func (d Dialect) Or(a ...Builder) *List {
	return Or(a...).Placeholder(d.Formatter)
}

func (d Dialect) OrLine(a ...Builder) *List {
	return OrLine(a...).Placeholder(d.Formatter)
}

func (d Dialect) OrEq(name string, values ...any) *List {
	return OrEq(name, values...).Placeholder(d.Formatter)
}

var (
	// Question 占位符作为 (?)
	Question = Dialect{Formatter: questionFormat{}}

	// Dollar 占位符为 ($1, $2, $3)
	Dollar = Dialect{Formatter: dollarFormat{}}

	// Colon 占位符为 ($1, $2, $3）
	Colon = Dialect{Formatter: colonFormat{}}

	// AtP 占位符为 (@p1, @p2, @p3)
	AtP = Dialect{Formatter: atpFormat{}}
)

type questionFormat struct{}

func (questionFormat) Placeholder() string {
	return "?"
}

func (questionFormat) ReplacePlaceholders(sql string) string {
	return sql
}

type dollarFormat struct{}

func (dollarFormat) Placeholder() string {
	return "$"
}

func (dollarFormat) ReplacePlaceholders(sql string) string {
	return replacePositionalPlaceholders(sql, "$")
}

type colonFormat struct{}

func (colonFormat) Placeholder() string {
	return ":"
}

func (colonFormat) ReplacePlaceholders(sql string) string {
	return replacePositionalPlaceholders(sql, ":")
}

type atpFormat struct{}

func (atpFormat) Placeholder() string {
	return "@p"
}

func (atpFormat) ReplacePlaceholders(sql string) string {
	return replacePositionalPlaceholders(sql, "@p")
}

func replacePositionalPlaceholders(sql, prefix string) string {
	buf := &bytes.Buffer{}
	i := 0
	for {
		p := strings.Index(sql, "?")
		if p == -1 {
			break
		}

		if len(sql[p:]) > 1 && sql[p:p+2] == "??" { // escape ?? => ?
			buf.WriteString(sql[:p])
			buf.WriteString("?")
			if len(sql[p:]) == 1 {
				break
			}
			sql = sql[p+2:]
		} else {
			i++
			buf.WriteString(sql[:p])
			_, _ = fmt.Fprintf(buf, "%s%d", prefix, i)
			sql = sql[p+1:]
		}
	}

	buf.WriteString(sql)
	return buf.String()
}
