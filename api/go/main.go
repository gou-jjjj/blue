package main

import "fmt"

func main() {
	c := NewClient(WithDefaultOpt())
	fmt.Println(c.Select())
	fmt.Println(c.Select("43"))
	fmt.Println(c.Select("3"))
	fmt.Println(c.Select())
	fmt.Println(c.Select("43"))
	fmt.Println(c.Select("3"))
}
