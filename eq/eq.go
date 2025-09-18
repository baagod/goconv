package eq

import (
	"reflect"
)

// Eq 等于 =
func Eq[T any](name string, value T, skip ...func(v T) bool) *Cond[T] {
	return NewCond(name, value, "=", skip...)
}

// Gt 大于 >
func Gt[T any](name string, value T, skip ...func(v T) bool) *Cond[T] {
	return NewCond(name, value, ">", skip...)
}

// Ge 大于等于 ≥
func Ge[T any](name string, value T, skip ...func(v T) bool) *Cond[T] {
	return NewCond(name, value, ">=", skip...)
}

// Lt 小于 <
func Lt[T any](name string, value T, skip ...func(v T) bool) *Cond[T] {
	return NewCond(name, value, "<", skip...)
}

// Le 小于等于 ≤
func Le[T any](name string, value T, skip ...func(v T) bool) *Cond[T] {
	return NewCond(name, value, "<=", skip...)
}

// Ne 不等于 != <>
func Ne[T any](name string, value T, skip ...func(v T) bool) *Cond[T] {
	return NewCond(name, value, "<>", skip...)
}

// Between 区间
func Between[T any](name string, first, second T, skip ...func(first, second T) bool) *Cond[*BETWEEN] {
	isSkip := len(skip) > 0 && skip[0] != nil && skip[0](first, second)
	return NewCond(name, &BETWEEN{first, second}, "BETWEEN").Skip(isSkip)
}

func NotBetween[T any](name string, first, second T, skip ...func(first, second T) bool) *Cond[*BETWEEN] {
	isSkip := len(skip) > 0 && skip[0] != nil && skip[0](first, second)
	return NewCond(name, &BETWEEN{first, second}, "NOT BETWEEN").Skip(isSkip)
}

// In 集合
func In[T any](name string, in []T, skip ...func(in []T) bool) *Cond[*IN] {
	values := make([]any, len(in))
	for i := range in {
		values[i] = in[i]
	}
	isSkip := len(skip) > 0 && skip[0] != nil && skip[0](in)
	return NewCond(name, &IN{values: values}, "IN").Skip(isSkip)
}

func NotIn[T any](name string, in []T, skip ...func(in []T) bool) *Cond[*IN] {
	values := make([]any, len(in))
	for i := range in {
		values[i] = in[i]
	}
	isSkip := len(skip) > 0 && skip[0] != nil && skip[0](in)
	return NewCond(name, &IN{values: values}, "NOT IN").Skip(isSkip)
}

// Like 模糊查询
func Like(name, value string, skip ...func(v string) bool) *Cond[*LIKE] {
	isSkip := len(skip) > 0 && skip[0] != nil && skip[0](value)
	return NewCond(name, &LIKE{value}, "LIKE").Skip(isSkip)
}

func NotLike(name, value string, skip ...func(v string) bool) *Cond[*LIKE] {
	isSkip := len(skip) > 0 && skip[0] != nil && skip[0](value)
	return NewCond(name, &LIKE{value}, "NOT LIKE").Skip(isSkip)
}

func IsNull(name string, skip ...func() bool) *Cond[any] {
	return &Cond[any]{
		Name:     name,
		Operator: "IS NULL",
		IsSkip:   len(skip) > 0 && skip[0] != nil && skip[0](),
	}
}

func IsNotNull(name string, skip ...func() bool) *Cond[any] {
	return &Cond[any]{
		Name:     name,
		Operator: "IS NOT NULL",
		IsSkip:   len(skip) > 0 && skip[0] != nil && skip[0](),
	}
}

// Or 或者
func Or(a ...Builder) *List {
	return &List{Builders: a, Sep: "OR", placeholder: DefaultPlaceholder}
}

func OrLine(a ...Builder) *List {
	return &List{Builders: a, Sep: "OR", enter: true, placeholder: DefaultPlaceholder}
}

// OrEq 生成 "col = ? OR col = ? OR ..." 的条件。
// 这是一个为常见的 OR 等值匹配提供的便利函数。
func OrEq(name string, values ...any) *List {
	builders := make([]Builder, len(values))
	for i, v := range values {
		builders[i] = Eq(name, v)
	}
	return Or(builders...)
}

// And 并且
func And(a ...Builder) *List {
	return &List{Builders: a, Sep: "AND", placeholder: DefaultPlaceholder}
}

func Where(a ...Builder) *List {
	return &List{Builders: a, Sep: "AND", indent: 2, placeholder: DefaultPlaceholder}
}

func isZero(v any) bool {
	//goland:noinspection GoDfaConstantCondition
	if v == nil || v == false || v == 0 || v == "" {
		return true
	}
	return reflect.ValueOf(v).IsZero()
}

func isNil(v any) bool {
	return v == nil
}
