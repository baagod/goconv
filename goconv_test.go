package goconv

import (
	"fmt"
	"testing"

	"github.com/baagod/goconv/eq"
)

func TestName(t *testing.T) {
	cond := eq.Dollar.Where(
		eq.Eq("name", "NAME"),
		eq.Eq("age", 18),
	)

	sql, args := cond.SQL()
	fmt.Println(sql, args)
	// fmt.Println("--------------------")
	// fmt.Println(eq.Debug(cond))
	// (a = $1 AND a = $2) OR a = $3
}
