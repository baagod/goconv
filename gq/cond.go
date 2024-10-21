package gq

import (
	"fmt"
	"reflect"
	"strings"
)

type Cond struct {
	Name     string
	Value    any
	Operator string
	isSkip   bool
}

func (c *Cond) Skip(skip bool) *Cond {
	c.isSkip = skip
	return c
}

func (c *Cond) SkipZero() *Cond {
	c.isSkip = reflect.ValueOf(c.Value).IsZero()
	return c
}

func (c *Cond) SkipFn(f func(v any) bool) *Cond {
	c.isSkip = f(c.Value)
	return c
}

func (c *Cond) SQL() string {
	if c.IsSkip() {
		return ""
	}
	return fmt.Sprintf("%s %s %s", c.Name, c.Operator, c.StrValue())
}

func (c *Cond) IsSkip() bool {
	return c.isSkip || (c.Operator == "LINK" && c.Value == "")
}

func (c *Cond) StrValue() string {
	s := fmt.Sprintf("%v", c.Value)

	switch c.Operator {
	case "BETWEEN":
		a := c.Value.([2]any)
		return fmt.Sprintf("%v BETWEEN %v", a[0], a[1])
	case "IN":
		ss := strings.Split(strings.Trim(s, "[]"), " ")
		return "(" + strings.Join(ss, ", ") + ")"
	case "LIKE":
		if !strings.HasPrefix(s, "%") && !strings.HasSuffix(s, "%") {
			s = fmt.Sprint(`'%`, s, `%'`)
		}
	}

	if _, ok := c.Value.(string); ok {
		s = fmt.Sprint(`'`, s, `'`)
	}

	return s
}
