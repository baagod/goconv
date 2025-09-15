package eq

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Builder interface {
	SQL() (string, []any)
}

type Cond struct {
	Name     string
	Value    any
	Operator string
	IsSkip   bool
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

type IN struct {
	values []any
}

func (in *IN) SQL() (string, []any) {
	length := len(in.values)
	if length == 0 {
		return "", nil
	}

	placeholders := make([]string, length)
	for i := range in.values {
		placeholders[i] = "?"
	}

	return "(" + strings.Join(placeholders, ", ") + ")", in.values
}

type LIKE struct {
	value string
}

func (l *LIKE) SQL() (placeholder string, args []any) {
	value := l.value
	if !strings.HasPrefix(l.value, "%") && !strings.HasSuffix(l.value, "%") {
		value = "%" + l.value + "%"
	}
	return "?", []any{value}
}

func NewCond(name string, value any, operator string, skip ...bool) *Cond {
	return &Cond{Name: name, Value: value, Operator: operator, IsSkip: skip != nil && skip[0]}
}

func (c *Cond) Skip(skip bool) *Cond {
	c.IsSkip = skip
	return c
}

func (c *Cond) SkipZero() *Cond {
	c.IsSkip = c.Value == nil || reflect.ValueOf(c.Value).IsZero()
	return c
}

func (c *Cond) SkipFn(f func(v any) bool) *Cond {
	c.IsSkip = f(c.Value)
	return c
}

func (c *Cond) SQL() (sql string, args []any) {
	if c.IsSkip {
		return "", nil
	}

	if c.Operator == "IS NULL" || c.Operator == "IS NOT NULL" {
		return c.Name + " " + c.Operator, nil
	}

	if c.Value == nil {
		return "", nil
	}

	switch v := c.Value.(type) {
	case Builder:
		sql, args = v.SQL()
	default:
		sql, args = "?", []any{c.Value}
	}

	return fmt.Sprintf("%s %s %s", c.Name, c.Operator, sql), args
}

func (c *Cond) Debug() string {
	return debug(c.SQL())
}

// ----

type List struct {
	Builders []Builder
	Sep      string // WHERE 条件分隔符
	indent   int    // 用于存储每行开头的基础缩进
}

func (lst *List) Append(a ...Builder) *List {
	lst.Builders = append(lst.Builders, a...)
	return lst
}

func (lst *List) Indent(width int) *List {
	if width > 0 {
		lst.indent = width
	}
	return lst
}

func (lst *List) IsIndent() bool {
	return lst.indent > 0
}

func (lst *List) SQL() (string, []any) {
	if lst.IsIndent() {
		return lst.format()
	}

	var sqls []string
	var allArgs []any

	for _, x := range lst.Builders {
		if sql, args := x.SQL(); sql != "" {
			allArgs = append(allArgs, args...)
			sqls = append(sqls, sql)
		}
	}

	if len(sqls) == 0 {
		return "", nil
	}

	sql := strings.Join(sqls, fmt.Sprintf(" %s ", lst.Sep))
	if lst.Sep == "OR" && len(sqls) > 1 {
		sql = "(" + sql + ")"
	}

	return sql, allArgs
}

func (lst *List) Debug() string {
	if !lst.IsIndent() {
		return debug(lst.SQL())
	}
	return debug(lst.format())
}

func (lst *List) format(baseIndent ...int) (string, []any) {
	indentWidth := lst.indent
	if len(baseIndent) > 0 {
		indentWidth = baseIndent[0]
	}

	var allArgs []any
	var groups []*List
	parts := map[int]bool{}

	for i, x := range lst.Builders {
		sql, args := x.SQL()
		if sql == "" {
			continue
		}

		allArgs = append(allArgs, args...)
		if part, ok := x.(*List); ok && part.Sep == "AND" {
			parts[i] = true
		}

		if len(groups) == 0 || parts[i] || parts[i-1] {
			groups = append(groups, &List{Builders: []Builder{x}, Sep: lst.Sep})
		} else {
			groups[len(groups)-1].Append(x)
		}
	}

	var lines []string
	for _, group := range groups {
		var groupSQLs []string
		for _, b := range group.Builders {
			var sql string
			if x, ok := b.(*List); ok {
				sql, _ = x.format(indentWidth)
			} else {
				sql, _ = b.SQL()
			}

			groupSQLs = append(groupSQLs, sql)
		}

		line := strings.Join(groupSQLs, fmt.Sprintf(" %s ", group.Sep))
		lines = append(lines, line)
	}

	// 这里是为了对齐 AND，就是不知道有没有更好的办法。
	sep := strings.Repeat(" ", indentWidth)
	if lst.Sep == "OR" {
		sep += " "
	}
	sep = fmt.Sprintf("\n%s%s ", sep, lst.Sep)

	sql := strings.Join(lines, sep)
	if lst.Sep == "OR" && len(lst.Builders) > 1 {
		sql = "(" + sql + ")"
	}

	return sql, allArgs
}

func debug(sql string, args []any) string {
	// 1. 按占位符分割 SQL 模板
	parts := strings.Split(sql, "?")

	// 2. 安全检查：分割后的片段数必须比参数数量多 1
	//    例如: "a = ? AND b = ?" -> parts: ["a = ", " AND b = ", ""], len=3; args: [val1, val2], len=2
	if len(parts)-1 != len(args) {
		return fmt.Sprintf(
			"/* DEBUGGER WARNING: Mismatch between %d placeholders and %d arguments */ %s",
			len(parts)-1, len(args), sql,
		)
	}

	// 3. 使用 strings.Builder 高效地交替拼接
	var sb strings.Builder
	for i, x := range parts {
		if sb.WriteString(x); i < len(args) {
			sb.WriteString(value(args[i]))
		}
	}

	return sb.String()
}

func value(value any) string {
	if value == nil {
		return "NULL"
	}

	switch v := value.(type) {
	case string:
		return "'" + strings.ReplaceAll(v, "'", "''") + "'"
	case time.Time:
		return v.Format("'2006-01-02 15:04:05'")
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		return fmt.Sprint(v)
	default:
		return "'" + strings.ReplaceAll(fmt.Sprint(v), "'", "''") + "'"
	}
}
