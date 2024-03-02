package cluster

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestNewCluster(t *testing.T) {
	cluster := NewCluster("127.0.0.1", 8080, 3, 10)
	dial, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	dial.Write([]byte("+192.22.43.11:9000"))
	dial.Close()

	dial, err = net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	dial.Write([]byte("+192.22.43.11:9010"))
	dial.Close()

	dial, err = net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	dial.Write([]byte("-192.22.43.11:9010"))
	dial.Close()
	<-time.NewTimer(5 * time.Second).C
	fmt.Println(cluster.observers)
}
