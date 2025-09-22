package eq

import (
	"fmt"
	"testing"
)

func TestEq(t *testing.T) {
	where := And( // 只搜医生
		Like("d.name", ""),
		// Lt("d.id", 0).OmitZero(),
	)

	// where.Append(
	// 	Like("h.name", ""),
	// )

	sql, args := where.SQL()
	fmt.Println(sql, args)
	// fmt.Println("--------------------")
	// fmt.Println(Debug(cond))
	// (a = $1 AND a = $2) OR a = $3
}
