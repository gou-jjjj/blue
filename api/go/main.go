package main

import "fmt"

func main() {
	c := NewClient(WithDefaultOpt())
	fmt.Println(c.Version())
	fmt.Println(c.Select("111"))
	fmt.Println(c.Select())
}
