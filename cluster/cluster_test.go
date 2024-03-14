package cluster

import (
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestNewCluster(t *testing.T) {
	cluster := NewCluster("127.0.0.1", 8080, 3, 10)
	dial, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	fmt.Println(dial.Write([]byte("+192.22.43.11:9015\n")))
	fmt.Println(dial.Write([]byte("+192.22.43.11:9013\n")))
	fmt.Println(dial.Write([]byte("+192.22.43.11:9019\n")))
	fmt.Println(dial.Close())

	<-time.NewTimer(2 * time.Second).C
	
	m := map[string]int{}
	for i := 0; i < int(1e5); i++ {
		m[cluster.c.Get(strconv.Itoa(i))]++
	}
	fmt.Println(m)

	fmt.Println(cluster.observers)
}
