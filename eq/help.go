package eq

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Debug 将 Builder 转换为可直接执行的 SQL 调试字符串。
// 它会将 SQL 中的 "?" 占位符替换为实际参数值（并进行转义）。
// 要求：传入的 Builder 或其子元素，其 SQL() 方法必须返回带 "?" 占位符的 SQL。
func Debug(b Builder) string {
	var sql string
	var args []any

	if l, ok := b.(*List); ok {
		sql, args = l.format()
	} else {
		sql, args = b.SQL()
	}

	var sb strings.Builder
	sqls := strings.Split(sql, "?")

	for i, x := range sqls[:len(sqls)-1] {
		sb.WriteString(x)
		sb.WriteString(debug(args[i]))
	}

	sb.WriteString(sqls[len(sqls)-1])
	return sb.String()
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
	case []byte:
		return string(v)
	default:
		return "'" + strings.ReplaceAll(fmt.Sprint(v), "'", "''") + "'"
	}
}
