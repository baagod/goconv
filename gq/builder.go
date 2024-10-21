package gq

import "strings"

type Builder interface {
	SQL() string
}

type WhereBuilder struct {
	List []Builder
	Sep  string // WHERE 条件分隔符
}

func (c *WhereBuilder) SQL() (sql string) {
	sep := " " + c.Sep + " "
	for _, x := range c.List {
		if s := x.SQL(); s != "" {
			sql += s + sep
		}
	}

	if c.Sep == "OR" {
		sql = "(" + sql + ")"
	}

	return strings.TrimSuffix(sql, sep)
}
