package goconv

import (
    "fmt"
    "testing"

    "github.com/baagod/goconv/eq"
)

func TestName(t *testing.T) {
    cond := eq.Where(
        eq.Eq("name", "O'Malley"),
        eq.Eq("age1", 18),
        eq.OrGroup(
            eq.And(eq.Between("time", "2010", "2019")),
            // eq.Between("time", "2010", "2019"),
            eq.Eq("height", 165),
            eq.Eq("height", 175),
        ),
        eq.And(
            eq.In("children", 1, 2, 3),
            eq.Like("content", "你好吗？"),
        ),
        eq.Eq("age3", 18),
    )

    sql := cond.Debug()
    fmt.Println(sql)

    // format := cond.Format()
    // fmt.Println(format)
    /* 输出：
       WHERE name = 'O''Malley'
          OR time BETWEEN '2010' AND '2019'
         AND (height = 165 OR height = 175)
         AND children IN (1, 2, 3) AND content LIKE '%你好吗？%',
         AND age = 18
    */
}
