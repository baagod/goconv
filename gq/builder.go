package gq

import (
	"regexp"
	"strings"
)

type Builder interface {
	SQL() string
}

type WhereBuilder struct {
	Builders []Builder
}

func (c *WhereBuilder) SQL() (sql string) {
	re := regexp.MustCompile("^\n? +AND |\n? +AND $")
	for _, c := range c.Builders {
		s := c.SQL()
		if s == "" {
			continue
		}

		if _, ok := c.(*List); ok {
			s += "\n"
		}

		if strings.HasSuffix(s, "\n") {
			sql = re.ReplaceAllString(sql, "") + "\n  AND " + s + "  AND "
		} else {
			sql += s + " AND "
		}
	}

	return re.ReplaceAllString(sql, "")
}
