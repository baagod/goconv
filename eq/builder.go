package eq

import (
	"fmt"
	"strings"
)

// Builder 是 SQL 条件构造器的核心接口。
// 约定：任何直接实现此接口的类型，其 SQL() 方法必须返回带有 "?" 占位符的原始 SQL。
// 如需生成其他占位符（如 $1, :1），请使用 List 包装器（如 Dollar.Where(), Colon.And()）。
type Builder interface {
	SQL() (string, []any)
}

type Cond[T any] struct {
	Name     string
	Value    T
	Operator string
	IsOmit   bool
}

func NewCond[T any](name string, value T, operator string, skip ...func(v T) bool) *Cond[T] {
	return &Cond[T]{
		Name:     name,
		Value:    value,
		Operator: operator,
		IsOmit:   len(skip) > 0 && skip[0](value),
	}
}

func (c *Cond[T]) Omit(omit bool) *Cond[T] {
	c.IsOmit = omit
	return c
}

func (c *Cond[T]) OmitZero() *Cond[T] {
	c.IsOmit = isZero(c.Value)
	return c
}

func (c *Cond[T]) OmitFn(f func(value T) bool) *Cond[T] {
	if f != nil {
		c.IsOmit = f(c.Value)
	}
	return c
}

func (c *Cond[T]) SQL() (sql string, args []any) {
	if c.IsOmit {
		return "", nil
	}

	if c.Operator == "IS NULL" || c.Operator == "IS NOT NULL" {
		return c.Name + " " + c.Operator, nil
	}

	if isNil(c.Value) {
		return "", nil
	}

	switch v := any(c.Value).(type) {
	case Builder:
		sql, args = v.SQL()
	default:
		sql, args = c.toSQL()
	}

	return fmt.Sprintf("%s %s %s", c.Name, c.Operator, sql), args
}

func (c *Cond[T]) toSQL() (string, []any) {
	switch c.Operator {
	case "LIKE", "NOT LIKE":
		value := any(c.Value).(string)
		if !strings.HasPrefix(value, "%") && !strings.HasSuffix(value, "%") {
			value = "%" + value + "%"
		}
		return "?", []any{value}
	}

	return "?", []any{c.Value}
}

type BETWEEN struct {
	first  any
	second any
}

func (b *BETWEEN) SQL() (string, []any) {
	if b.first == nil || b.second == nil {
		return "", nil
	}
	return "? AND ?", []any{b.first, b.second}
}

type IN[T any] struct {
	values []T
}

func (in *IN[T]) SQL() (string, []any) {
	if len(in.values) == 0 {
		return "", nil
	}

	args := make([]any, len(in.values))
	placeholders := make([]string, len(in.values))

	for i, x := range in.values {
		placeholders[i] = "?"
		args[i] = x
	}

	return "(" + strings.Join(placeholders, ", ") + ")", args
}

// ----

type List struct {
	Builders []Builder
	Sep      string // WHERE 条件分隔符
	Dialect  Placeholder
	indent   int  // 用于存储每行开头的基础缩进
	enter    bool // 在格式化时，是否为该 List 开启新行
}

func (l *List) Placeholder(p Placeholder) *List {
	l.Dialect = p
	return l
}

func (l *List) Append(a ...Builder) *List {
	l.Builders = append(l.Builders, a...)
	return l
}

func (l *List) Indent(width int) *List {
	if width > 0 {
		l.indent = width
	}
	return l
}

func (l *List) Enter(enter bool) *List {
	l.enter = enter
	return l
}

func (l *List) IsIndent() bool {
	return l.indent > 0
}

func (l *List) SQL() (string, []any) {
	if l.IsIndent() {
		sql, args := l.format()
		return l.Dialect.ReplacePlaceholders(sql), args
	}

	var sqls []string
	var allArgs []any

	for _, x := range l.Builders {
		sql, args := x.SQL()
		if sql == "" {
			continue
		}

		sqls = append(sqls, sql)
		allArgs = append(allArgs, args...)
	}

	if len(sqls) == 0 {
		return "", nil
	}

	sql := strings.Join(sqls, fmt.Sprintf(" %s ", l.Sep))
	if l.Sep == "OR" && len(sqls) > 1 {
		sql = "(" + sql + ")"
	}

	return l.Dialect.ReplacePlaceholders(sql), allArgs
}

func (l *List) format(baseIndent ...int) (string, []any) {
	if len(l.Builders) == 0 {
		return "", nil
	}

	indent := l.indent
	if indent == 0 {
		if len(baseIndent) == 0 {
			baseIndent = append(baseIndent, 2)
		}
		indent = baseIndent[0]
	}

	var groups []*List
	var allArgs []any
	var lines []string
	breaks := map[int]bool{}

	for i, x := range l.Builders {
		group, ok := x.(*List)
		if ok && (group.Sep == "AND" || group.enter) {
			if len(group.Builders) == 0 {
				continue
			}
			breaks[i] = true
		}

		if len(groups) == 0 || breaks[i] || breaks[i-1] {
			groups = append(groups, &List{Builders: []Builder{x}, Sep: l.Sep})
		} else {
			groups[len(groups)-1].Append(x)
		}
	}

	for _, group := range groups {
		var sqls []string
		var args []any
		var sql string

		for _, b := range group.Builders {
			if x, ok := b.(*List); ok {
				sql, args = x.format(indent)
			} else {
				sql, args = b.SQL()
			}
			if sql != "" {
				allArgs = append(allArgs, args...)
				sqls = append(sqls, sql)
			}
		}

		if len(sqls) > 0 {
			line := strings.Join(sqls, fmt.Sprintf(" %s ", group.Sep))
			lines = append(lines, line)
		}
	}

	if len(lines) == 0 {
		return "", nil
	}

	// 这里是为了对齐 AND，就是不知道有没有更好的办法。
	sep := strings.Repeat(" ", indent)
	if l.Sep == "OR" {
		sep += " "
	}
	sep = fmt.Sprintf("\n%s%s ", sep, l.Sep)

	sql := strings.Join(lines, sep)
	if l.Sep == "OR" && len(l.Builders) > 1 {
		sql = "(" + sql + ")"
	}

	return sql, allArgs
}
