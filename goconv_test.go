package goconv

import (
	"fmt"
	"testing"

	"github.com/baagod/goconv/gq"
)

func TestName(t *testing.T) {
	a := gq.Like("name", "sss")
	fmt.Println(a.SQL())
}
