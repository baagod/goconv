package goconv

import (
	"fmt"
	"testing"

	"github.com/baagod/goconv/eq"
)

func TestName(t *testing.T) {
	cond := eq.Dollar.Where(
		// eq.Eq("name", "NAME"),
		// eq.Eq("age", 18),
		eq.OrLine(
			eq.And(
				eq.Eq("a", 1),
				eq.Eq("a", 2),
			),
			eq.Eq("a", 2),
		),
		// eq.And(
		// 	eq.Eq("a", 2),
		// 	eq.Eq("a", 2),
		// 	eq.And(eq.Eq("b", 3)),
		// ),
	)

	sql, args := cond.SQL()
	fmt.Println(sql, args)
	// fmt.Println("--------------------")
	// fmt.Println(eq.Debug(cond))
	// (a = $1 AND a = $2) OR a = $3
}
