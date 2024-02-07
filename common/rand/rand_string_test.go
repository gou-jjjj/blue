package rand

import (
	"fmt"
	"testing"
)

func TestRandString(t *testing.T) {
	for i := 0; i < 32; i++ {
		randString := RandString(i)
		fmt.Printf("len: %d | randString: %s\n", len(randString), randString)
	}
}
