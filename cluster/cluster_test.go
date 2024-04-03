package cluster

import (
	"fmt"
	"testing"
	"time"
)

func TestNewCluster(t *testing.T) {
	cluster := NewCluster(1, 8080, "1", 1)

	cluster.Register(
		"127.0.0.1:8080",
		"127.0.0.1:8083")

	fmt.Println(cluster.observers)
	c := NewCluster(1, 8081, "1", 1)

	c.GetClusterAddrs("10.16.101.0:8080")

	cluster.Register("aa")
	c.GetClusterAddrs("10.16.101.0:8080")

	cluster.Register("aa1")

	time.Sleep(4 * time.Second)
}
