package internal

import (
	"blue/common/network"
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Printf("%v\n", network.LocalIpEn0())
}
