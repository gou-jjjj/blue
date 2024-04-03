package cluster

import (
	"fmt"
	"testing"
	"time"
)

func TestNewCluster(t *testing.T) {

	c := NewCluster(1, 8081, "1", 1)
	fmt.Println(c.observers)
	//c.GetClusterAddrs("39.101.195.49:7891")

	c.Register("a", "b")
	c.Register("a", "b")
	c.Register("a", "b")
	c.Register("a", "b", "c")
	c.Register("a", "b")
	c.Register("f")

	fmt.Println(c.observers)

	c.Unregister("a", "b")
	c.Unregister("a", "f")
	c.Unregister("a", "b")
	fmt.Println(c.observers)
	time.Sleep(4 * time.Second)
}
