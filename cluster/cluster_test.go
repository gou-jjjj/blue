package cluster

import (
	"fmt"
	"testing"
	"time"
)

func TestGetAddrs(t *testing.T) {
	c := NewCluster(1, "127.0.0.1", 8081, "127.0.0.1:8080", "127.0.0.1:8082", 1*time.Minute)

	addrs := c.GetClusterAddr("127.0.0.1:7892")

	c.Register(addrs...)
	fmt.Printf("%+v\n", c.observers)

}
