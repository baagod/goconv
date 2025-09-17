package eq

import (
	"reflect"
)

// Eq 等于 =
func Eq[T any](name string, value T, skip ...func(v T) bool) *Cond {
	return NewCond(name, value, "=", skip...)
}

// Gt 大于 >
func Gt[T any](name string, value T, skip ...func(v T) bool) *Cond {
	return NewCond(name, value, ">", skip...)
}

// Ge 大于等于 ≥
func Ge[T any](name string, value T, skip ...func(v T) bool) *Cond {
	return NewCond(name, value, ">=", skip...)
}

// Lt 小于 <
func Lt[T any](name string, value T, skip ...func(v T) bool) *Cond {
	return NewCond(name, value, "<", skip...)
}

// Le 小于等于 ≤
func Le[T any](name string, value T, skip ...func(v T) bool) *Cond {
	return NewCond(name, value, "<=", skip...)
}

// Ne 不等于 != <>
func Ne[T any](name string, value T, skip ...func(v T) bool) *Cond {
	return NewCond(name, value, "<>", skip...)
}

// Between 区间
func Between[F any, S any](name string, first F, second S, skip ...func(first F, second S) bool) *Cond {
	isSkip := isNil(first) || isNil(second)
	if !isSkip && len(skip) > 0 && skip[0] != nil {
		isSkip = skip[0](first, second)
	}
	return NewCond(name, &BETWEEN{first, second}, "BETWEEN").Skip(isSkip)
}

func NotBetween[F any, S any](name string, first F, second S, skip ...func(first F, second S) bool) *Cond {
	isSkip := isNil(first) || isNil(second)
	if !isSkip && len(skip) > 0 && skip[0] != nil {
		isSkip = skip[0](first, second)
	}
	return NewCond(name, &BETWEEN{first, second}, "NOT BETWEEN").Skip(isSkip)
}

// In 集合
func In[T any](name string, in []T, skip ...func(in []T) bool) *Cond {
	isSkip := len(in) == 0
	if !isSkip && len(skip) > 0 && skip[0] != nil {
		isSkip = skip[0](in)
	}

	values := make([]any, len(in))
	for i, v := range in {
		values[i] = v
	}

	return NewCond(name, &IN{values: values}, "IN").Skip(isSkip)
}

func NotIn[T any](name string, in []T, skip ...func(in []T) bool) *Cond {
	isSkip := len(in) == 0
	if !isSkip && len(skip) > 0 && skip[0] != nil {
		isSkip = skip[0](in)
	}

	values := make([]any, len(in))
	for i, v := range in {
		values[i] = v
	}

	return NewCond(name, &IN{values: values}, "NOT IN").Skip(isSkip)
}

// Like 模糊查询
func Like(name, value string, skip ...func(v string) bool) *Cond {
	isSkip := value == ""
	if !isSkip && len(skip) > 0 && skip[0] != nil {
		isSkip = skip[0](value)
	}
	return NewCond(name, &LIKE{value}, "LIKE").Skip(isSkip)
}

func NotLike(name, value string, skip ...func(v string) bool) *Cond {
	isSkip := value == ""
	if !isSkip && len(skip) > 0 && skip[0] != nil {
		isSkip = skip[0](value)
	}
	return NewCond(name, &LIKE{value}, "NOT LIKE").Skip(isSkip)
}

func IsNull(name string, skip ...func() bool) *Cond {
	return &Cond{
		Name:     name,
		Operator: "IS NULL",
		IsSkip:   len(skip) == 0 || skip[0] == nil || skip[0](),
	}
}

func IsNotNull(name string, skip ...func() bool) *Cond {
	return &Cond{
		Name:     name,
		Operator: "IS NOT NULL",
		IsSkip:   len(skip) == 0 || skip[0] == nil || skip[0](),
	}
}

// Or 或者
func Or(a ...Builder) *List {
	return &List{Builders: a, Sep: "OR", placeholder: questionFormat{}}
}

func OrLine(a ...Builder) *List {
	return &List{Builders: a, Sep: "OR", isGroup: true, placeholder: questionFormat{}}
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
	return &List{Builders: a, Sep: "AND", placeholder: questionFormat{}}
}

func Where(a ...Builder) *List {
	return &List{Builders: a, Sep: "AND", indent: 2, placeholder: questionFormat{}}
}

func isZero(v any) bool {
	//goland:noinspection GoDfaConstantCondition
	if v == nil || v == false || v == 0 || v == "" {
		return true
	}
	return reflect.ValueOf(v).IsZero()
}

func isNil(v any) bool {
	if v == nil { // 处理真正的 nil 接口
		return true
	}
	switch val := reflect.ValueOf(v); val.Kind() { // 处理像 (*int)(nil) 这种情况
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr,
		reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
}
