package main

import "fmt"

func main() {
	c := NewClient(WithDefaultOpt())
	fmt.Println(c.Version())
	fmt.Println(c.Select("1"))
	fmt.Println(c.Select())
	fmt.Println(c.Select("2"))
	fmt.Println(c.Select())
}
