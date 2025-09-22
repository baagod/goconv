package eq

import (
	"fmt"
	"testing"
)

func TestEq(t *testing.T) {
	cond := Dollar.Where(
		Eq("name", "NAME"),
		Eq("age", 18),
		In("a", []int{3, 4}),
	)

	sql, args := cond.SQL()
	fmt.Println(sql, args)
	// fmt.Println("--------------------")
	// fmt.Println(Debug(cond))
	// (a = $1 AND a = $2) OR a = $3
}
