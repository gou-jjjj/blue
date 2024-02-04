package main

import (
	"fmt"
)

func RedMessage(s string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", s)
}

func GreenMessage(s string) string {
	return fmt.Sprintf("\033[32m%s\033[0m", s)
}

func YellowMessage(s string) string {
	return fmt.Sprintf("\033[33m%s\033[0m", s)
}

func BlueMessage(s string) string {
	return fmt.Sprintf("\033[34m%s\033[0m", s)
}

func ErrPrint(d any) {
	fmt.Printf("%v\n", d)
}

func SuccessPrint(d any) {
	fmt.Printf("%v\n", d)
}
