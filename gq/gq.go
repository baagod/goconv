package gq

import (
	"fmt"
	"strings"
)

// Eq 等于 =
func Eq[T any](name string, v T, skip ...func(T) bool) *Cond {
	return &Cond{Name: name, Value: v, Opr: "=", isSkip: skip != nil && skip[0](v)}
}

// Gt 大于 >
func Gt[T any](name string, v T, skip ...func(T) bool) *Cond {
	return &Cond{Name: name, Value: v, Opr: ">", isSkip: skip != nil && skip[0](v)}
}

// Ge 大于等于 ≥
func Ge[T any](name string, v T, skip ...func(T) bool) *Cond {
	return &Cond{Name: name, Value: v, Opr: ">=", isSkip: skip != nil && skip[0](v)}
}

// Lt 小于 <
func Lt[T any](name string, v T, skip ...func(T) bool) *Cond {
	return &Cond{Name: name, Value: v, Opr: "<", isSkip: skip != nil && skip[0](v)}
}

// Le 小于等于 ≤
func Le[T any](name string, v T, skip ...func(T) bool) *Cond {
	return &Cond{Name: name, Value: v, Opr: "<=", isSkip: skip != nil && skip[0](v)}
}

// Ne 不等于 != <>
func Ne[T any](name string, v T, skip ...func(T) bool) *Cond {
	return &Cond{Name: name, Value: v, Opr: "!=", isSkip: skip != nil && skip[0](v)}
}

// Between 区间
func Between[T any](name string, first, second T, skip ...func(T, T) bool) *Cond {
	return &Cond{
		Name:   name,
		Value:  [2]any{first, second},
		Opr:    "BETWEEN",
		isSkip: skip != nil && skip[0](first, second),
	}
}

// In 范围
func In[T any](name string, a ...T) *Cond {
	if a != nil {
		// 是否切片或数组
		if strings.HasPrefix(fmt.Sprintf("%T", a[0]), "[") {
			return &Cond{Name: name, Value: a[0], Opr: "IN"}
		}
	}
	return &Cond{Name: name, Value: a, Opr: "IN"}
}

// Like 模糊查询
func Like(name string, v string, skip ...func(string) bool) *Cond {
	return &Cond{Name: name, Value: v, Opr: "LIKE", isSkip: skip != nil && skip[0](v)}
}

// Or 或者条件
func Or(a ...Builder) *List {
	return &List{Builders: a, Sep: "OR"}
}

// And 并且条件
func And(a ...Builder) *List {
	return &List{Builders: a, Sep: "AND"}
}

func Where(a ...Builder) *WhereBuilder {
	return &WhereBuilder{Builders: a}
}
