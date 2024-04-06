package testclu

import (
	blue "blue/api/go"
	"blue/cluster"
	"fmt"
	"strconv"
	"testing"
	"time"
)

var localclu1 = "127.0.0.1:7891"
var localclu2 = "127.0.0.1:7892"

var localclu3 = "127.0.0.1:7893"

func TestGetClu(t *testing.T) {
	blu, err := blue.NewClient(blue.WithDefaultOpt(), blue.WithAddr(localclu3))
	if err != nil {
		t.Error(err)
	}
	defer blu.Close()

	for i := 0; i < 10000; i++ {
		v, err := blu.Get(strconv.Itoa(i))
		if err != nil {
			t.Error(err)
		}

		if v != "b" {
			t.Error("bbb", v)
		}

	}
}

func TestJoinClu(t *testing.T) {
	c := cluster.NewCluster(1, "127.0.0.1", 8081, "127.0.0.1:8081", "127.0.0.1:8082", 10*time.Second)
	addrs1 := c.GetClusterAddr(localclu2)

	c.InitClusterAddr(addrs1...)

	fmt.Printf("%+v\n %+v\n", c.GetObservers(), c.GetConsistent().Members())

	c.GetClusterAddr(localclu1)
	c.GetClusterAddr(localclu2)

	c.Close()
}
