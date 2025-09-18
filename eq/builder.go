package eq

import (
    "fmt"
    "strings"
    "time"
)

type Builder interface {
    SQL() (string, []any)
}

type Cond[T any] struct {
    Name     string
    Value    T
    Operator string
    IsSkip   bool
}

func NewCond[T any](name string, value T, operator string, skip ...func(v T) bool) *Cond[T] {
    return &Cond[T]{
        Name:     name,
        Value:    value,
        Operator: operator,
        IsSkip:   len(skip) > 0 && skip[0] != nil && skip[0](value),
    }
}

func (c *Cond[T]) Skip(skip bool) *Cond[T] {
    c.IsSkip = skip
    return c
}

func (c *Cond[T]) SkipZero() *Cond[T] {
    c.IsSkip = isZero(c.Value)
    return c
}

func (c *Cond[T]) SkipFn(f func(value T) bool) *Cond[T] {
    if f != nil {
        c.IsSkip = f(c.Value)
    }
    return c
}

func (c *Cond[T]) SQL() (sql string, args []any) {
    if c.IsSkip {
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
        sql, args = "?", []any{c.Value}
    }

    return fmt.Sprintf("%s %s %s", c.Name, c.Operator, sql), args
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

func (l *LIKE) SQL() (string, []any) {
    if l.value == "" {
        return "", nil
    }

    value := l.value
    if !strings.HasPrefix(l.value, "%") && !strings.HasSuffix(l.value, "%") {
        value = "%" + l.value + "%"
    }

    return "?", []any{value}
}

// ----

type List struct {
    Builders    []Builder
    Sep         string // WHERE 条件分隔符
    indent      int    // 用于存储每行开头的基础缩进
    enter       bool   // 在格式化时，是否为该 List 开启一个新行
    placeholder Placeholder
}

func (l *List) Placeholder(format Placeholder) *List {
    l.placeholder = format
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

func (l *List) IsIndent() bool {
    return l.indent > 0
}

func (l *List) SQL() (string, []any) {
    sql, args := l.format()
    return l.placeholder.ReplacePlaceholders(sql), args
}

// collect 辅助方法
func (l *List) collect() (sqls []string, groups []*List, allArgs []any) {
    parts := map[int]bool{}
    isFormat := l.IsIndent()

    for i, x := range l.Builders {
        sql, args := x.SQL()
        if sql == "" {
            continue
        }

        if allArgs = append(allArgs, args...); !isFormat {
            sqls = append(sqls, sql)
            continue
        }

        if part, ok := x.(*List); ok && (part.Sep == "AND" || part.enter) {
            parts[i] = true
        }

        if len(groups) == 0 || parts[i] || parts[i-1] {
            groups = append(groups, &List{Builders: []Builder{x}, Sep: l.Sep})
        } else {
            groups[len(groups)-1].Append(x)
        }
    }

    return
}

func (l *List) format(baseIndent ...int) (string, []any) {
    sqls, groups, allArgs := l.collect()

    if len(groups) == 0 { // 非格式化 SQL
        if len(sqls) == 0 {
            return "", nil
        }

        sql := strings.Join(sqls, fmt.Sprintf(" %s ", l.Sep))
        if l.Sep == "OR" && len(sqls) > 1 {
            sql = "(" + sql + ")"
        }

        return sql, allArgs
    }

    indent := l.indent
    if len(baseIndent) > 0 {
        indent = baseIndent[0]
    }

    var lines []string
    for _, group := range groups {
        var groupSQLs []string
        for _, b := range group.Builders {
            var sql string
            if x, ok := b.(*List); ok {
                sql, _ = x.format(indent)
            } else {
                sql, _ = b.SQL()
            }

            groupSQLs = append(groupSQLs, sql)
        }

        line := strings.Join(groupSQLs, fmt.Sprintf(" %s ", group.Sep))
        lines = append(lines, line)
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

func Debug(b Builder) string {
    var sql string
    var args []any

    if l, ok := b.(*List); ok {
        sql, args = l.format()
    } else {
        sql, args = b.SQL()
    }

    // 安全检查：分割后的片段数必须比参数数量多 1
    // 例如: "a = ? AND b = ?" -> conds: ["a = ", " AND b = ", ""], len=3; args: [val1, val2], len=2
    conds := strings.Split(sql, "?")
    if len(conds)-1 != len(args) {
        return fmt.Sprintf(
            "/* DEBUGGER WARNING: Mismatch between %d placeholders and %d arguments */ %s",
            len(conds)-1, len(args), sql,
        )
    }

    // 使用 strings.Builder 高效地交替拼接
    var sb strings.Builder
    for i, x := range conds {
        if sb.WriteString(x); i < len(args) {
            sb.WriteString(debug(args[i]))
        }
    }

    return sb.String()
}

func debug(value any) string {
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
