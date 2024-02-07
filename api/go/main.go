package main

import "fmt"

func main() {
	c := NewClient(WithDefaultOpt())
	fmt.Println(c.Version())
}
