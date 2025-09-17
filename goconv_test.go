package goconv

import (
	"fmt"
	"testing"

	"github.com/baagod/goconv/eq"
)

func TestName(t *testing.T) {
	cond := eq.Dollar.And(
		eq.Eq("name", "NAME"),
		eq.Eq("age", 18),
		eq.And(
			eq.Eq("age", 12),
		),
	)

	sql, args := cond.SQL()
	fmt.Println(sql, args)
	fmt.Println("--------------------")
	fmt.Println(eq.Debug(cond))

	// squirrel.DebugSqlizer()
}
