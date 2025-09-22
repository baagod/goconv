package eq

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
	isSkip := len(skip) > 0 && skip[0](first, second)
	return NewCond(name, &BETWEEN{first, second}, "BETWEEN").Omit(isSkip)
}

func NotBetween[T any](name string, first, second T, skip ...func(first, second T) bool) *Cond[*BETWEEN] {
	isSkip := len(skip) > 0 && skip[0](first, second)
	return NewCond(name, &BETWEEN{first, second}, "NOT BETWEEN").Omit(isSkip)
}

// In 集合
func In[T any](name string, in []T, skip ...func(in []T) bool) *Cond[*IN[T]] {
	return NewCond(name, &IN[T]{in}, "IN").Omit(len(skip) > 0 && skip[0](in))
}

func NotIn[T any](name string, in []T, skip ...func(in []T) bool) *Cond[*IN[T]] {
	return NewCond(name, &IN[T]{in}, "NOT IN").Omit(len(skip) > 0 && skip[0](in))
}

// Like 模糊查询
func Like(name, value string, skip ...func(v string) bool) *Cond[string] {
	return NewCond(name, value, "LIKE", skip...)
}

func NotLike(name, value string, skip ...func(v string) bool) *Cond[string] {
	return NewCond(name, value, "NOT LIKE", skip...)
}

func IsNull(name string, skip ...func() bool) *Cond[any] {
	return &Cond[any]{
		Name:     name,
		Operator: "IS NULL",
		IsOmit:   len(skip) > 0 && skip[0](),
	}
}

func IsNotNull(name string, skip ...func() bool) *Cond[any] {
	return &Cond[any]{
		Name:     name,
		Operator: "IS NOT NULL",
		IsOmit:   len(skip) > 0 && skip[0](),
	}
}

// Or 或者
func Or(a ...Builder) *List {
	return &List{Builders: a, Sep: "OR", Dialect: DefaultPlaceholder}
}

func OrLine(a ...Builder) *List {
	return Or(a...).Enter(true)
}

// And 并且
func And(a ...Builder) *List {
	return &List{Builders: a, Sep: "AND", Dialect: DefaultPlaceholder}
}

func AndLine(a ...Builder) *List {
	return And(a...).Enter(true)
}

func Where(a ...Builder) *List {
	return And(a...).Indent(2)
}
