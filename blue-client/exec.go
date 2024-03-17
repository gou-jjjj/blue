package main

import (
	g "blue/api/go"
)

type CmdFunc func(*g.Client, []string) (string, error)

func Exec(conn *g.Client, s []string) (string, error) {
	if len(s) == 0 {
		return "", ErrCommand(s[0])
	}

	f, ok := funcMap[s[0]]
	if !ok {
		return "", ErrCommand(s[0])
	}
	return f(conn, s)
}

func set() CmdFunc {
	return func(conn *g.Client, s []string) (string, error) {
		if len(s) != 3 {
			return "", ErrArgu(s[0])
		}
		return conn.Set(s[1], s[2])
	}
}

var funcMap = map[string]CmdFunc{
	"set": set(),
}
