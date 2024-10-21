package gq

import (
	"fmt"
	"regexp"
	"strings"
)

// Eq 等于 =
func Eq(name string, v any) *Cond {
	return &Cond{Name: name, Value: v, Operator: "="}
}

// Gt 大于 >
func Gt(name string, v any) *Cond {
	return &Cond{Name: name, Value: v, Operator: ">"}
}

// Ge 大于等于 ≥
func Ge(name string, v any) *Cond {
	return &Cond{Name: name, Value: v, Operator: ">="}
}

// Lt 小于 <
func Lt(name string, v any) *Cond {
	return &Cond{Name: name, Value: v, Operator: "<"}
}

// Le 小于等于 ≤
func Le(name string, v any) *Cond {
	return &Cond{Name: name, Value: v, Operator: "<="}
}

// Ne 不等于 != <>
func Ne(name string, v any) *Cond {
	return &Cond{Name: name, Value: v, Operator: "!="}
}

// Between 区间
func Between(name string, first, second any) *Cond {
	return &Cond{
		Name:     name,
		Value:    [2]any{first, second},
		Operator: "BETWEEN",
	}
}

// In 范围
func In(name string, a ...any) *Cond {
	if len(a) == 1 {
		if strings.HasPrefix(fmt.Sprintf("%T", a[0]), "[") {
			return &Cond{Name: name, Value: a[0], Operator: "IN"}
		}
	}
	return &Cond{Name: name, Value: a, Operator: "IN"}
}

// Like 模糊查询
func Like(name string, v string) *Cond {
	return &Cond{Name: name, Value: v, Operator: "LIKE"}
}

// Or 或者条件
func Or(a ...Builder) *WhereBuilder {
	return &WhereBuilder{List: a, Sep: "OR"}
}

// And 并且条件
func And(a ...Builder) *WhereBuilder {
	return &WhereBuilder{List: a, Sep: "AND"}
}

func Where(a ...Builder) (sql string) {
	re := regexp.MustCompile("^\n? +AND |\n? +AND $")
	for _, c := range a {
		s := c.SQL()
		if s == "" {
			continue
		}

		if _, ok := c.(*WhereBuilder); ok {
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
