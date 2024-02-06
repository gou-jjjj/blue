package main

import (
	"fmt"
)

func ErrArgu(s ...string) error {
	return fmt.Errorf("(error)  wrong number of arguments for '%s' command", RedMessage(s[0]))
}

func ErrType(s ...string) error {
	return fmt.Errorf("(error)  unknown type '%s'", RedMessage(s[0]))
}

func ErrCommandType(s ...string) error {
	return fmt.Errorf("(error)  unknown command type '%s'", RedMessage(s[0]))
}

func ErrDataType(s ...string) error {
	return fmt.Errorf("(error)  unknown data type '%s'", RedMessage(s[0]))
}

func ErrCommand(s ...string) error {
	return fmt.Errorf("(error)  unknown command '%s'", RedMessage(s[0]))
}

func ErrSyntax(s ...string) error {
	return fmt.Errorf("(error)  syntax error '%s'", RedMessage(s[0]))
}

func ErrCommandNil(s ...string) error {
	return fmt.Errorf("(error)  command is '%s'", RedMessage("nil"))
}

func ErrConnect(s ...string) error {
	return fmt.Errorf("(error)  connect error '%s'", RedMessage(s[0]))
}

func ErrRead(s ...string) error {
	return fmt.Errorf("(error)  read command error '%s'", RedMessage(s[0]))
}

func ErrInvalidResp(s ...string) error {
	return fmt.Errorf("(error)  %s", RedMessage("invalid response"))
}
