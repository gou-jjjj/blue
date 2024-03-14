package main

import (
	g "blue/api/go"
)

type CmdFunc func([]string) (string, error)

var funcMap = map[string]CmdFunc{}

func Register(name string, f CmdFunc) {
	funcMap[name] = f
}

func set(conn g.Client) CmdFunc {
	return func(s []string) (string, error) {
		if len(s) != 3 {
			return "", ErrArgu(s[0])
		}
		return conn.Set(s[1], s[2])
	}
}

// Exec is a function that takes a slice of strings and returns a string and an error.
func Exec(s []string) (string, error) {
	if len(s) == 0 {
		return "", ErrCommand(s[0])
	}

	f, ok := funcMap[s[0]]
	if !ok {
		return "", ErrCommand(s[0])
	}
	return f(s)
}

// len(s) == 0 {
//	return "", ErrCommand(s[0])
//}
//
//switch s[0] {
//case "exit":
//	os.Exit(0)
//	//return "", nil
//case "set":
//	if len(s) != 3 {
//		return "", ErrArgu(s[0])
//	}
//	return conn.Set(s[1], s[2])
//case "del":
//	if len(s) != 2 {
//		return "", ErrArgu(s[0])
//	}
//	return conn.Del(s[1])
//case "version":
//	if len(s) != 1 {
//		return "", ErrArgu(s[0])
//	}
//	return conn.Version()
//case "get":
//	if len(s) != 2 {
//		return "", ErrArgu(s[0])
//	}
//	return conn.Get(s[1])
//
//default:
//	return "", ErrCommand(s[0])
//}
//return "", nil
