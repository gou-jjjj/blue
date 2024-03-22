package cluster

import (
	"fmt"
	"testing"
)

func TestNewCluster(t *testing.T) {
	cluster := NewCluster(1, 8080, "1", 1)

	fmt.Printf("%v\n", cluster.observers)
	fmt.Printf("%v\n", cluster.c.members)

	cluster.Register("dqwdq", "dwd")
	fmt.Printf("%v\n", cluster.observers)
	fmt.Printf("%v\n", cluster.c.members)

}
