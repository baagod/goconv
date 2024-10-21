package goconv

import (
	"fmt"
	"testing"

	"github.com/baagod/goconv/strmu"
)

func TestName(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(strmu.Rand(12, true))
	}
}
