package eq

// Eq 等于 =
func Eq(name string, value any) *Cond {
	return NewCond(name, value, "=")
}

// Gt 大于 >
func Gt(name string, value any) *Cond {
	return NewCond(name, value, ">")
}

// Ge 大于等于 ≥
func Ge(name string, value any) *Cond {
	return NewCond(name, value, ">=")
}

// Lt 小于 <
func Lt(name string, value any) *Cond {
	return NewCond(name, value, "<")
}

// Le 小于等于 ≤
func Le(name string, value any) *Cond {
	return NewCond(name, value, "<=")
}

// Ne 不等于 != <>
func Ne(name string, value any) *Cond {
	return NewCond(name, value, "<>")
}

// Between 区间
func Between(name string, first, second any) *Cond {
	value := &BETWEEN{first, second}
	return NewCond(name, value, "BETWEEN", first == nil || second == nil)
}

func NotBetween(name string, first, second any) *Cond {
	return NewCond(name, &BETWEEN{first, second}, "NOT BETWEEN", first == nil || second == nil)
}

// In 集合
func In(name string, a ...any) *Cond {
	return NewCond(name, &IN{a}, "IN", a == nil || len(a) == 0)
}

func NotIn(name string, a ...any) *Cond {
	return NewCond(name, &IN{a}, "NOT IN", a == nil || len(a) == 0)
}

// Like 模糊查询
func Like(name, value string) *Cond {
	return NewCond(name, &LIKE{value}, "LIKE", value == "")
}

func NotLike(name, value string) *Cond {
	return NewCond(name, &LIKE{value}, "NOT LIKE", value == "")
}

func IsNull(name string) *Cond {
	return &Cond{Name: name, Operator: "IS NULL"}
}

func IsNotNull(name string) *Cond {
	return &Cond{Name: name, Operator: "IS NOT NULL"}
}

// Or 或者条件
func Or(a ...Builder) *List {
	return &List{Builders: a, Sep: "OR"}
}

// And 并且条件
func And(a ...Builder) *List {
	return &List{Builders: a, Sep: "AND"}
}

func Where(a ...Builder) *List {
	return &List{Builders: a, Sep: "AND", indent: 2}
}
