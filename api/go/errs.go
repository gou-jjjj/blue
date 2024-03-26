package blue

import (
	"fmt"
)

func ErrArgu(s ...string) error {
	return fmt.Errorf("(error)  wrong number of arguments for '%s' command", s[0])
}

func ErrType(s ...string) error {
	return fmt.Errorf("(error)  unknown type '%s'", s[0])
}

func ErrCommandType(s ...string) error {
	return fmt.Errorf("(error)  unknown command type '%s'", (s[0]))
}

func ErrDataType(s ...string) error {
	return fmt.Errorf("(error)  unknown data type '%s'", (s[0]))
}

func ErrCommand(s ...string) error {
	return fmt.Errorf("(error)  unknown command '%s'", (s[0]))
}

func ErrSyntax(s ...string) error {
	return fmt.Errorf("(error)  syntax error '%s'", (s[0]))
}

func ErrCommandNil(s ...string) error {
	return fmt.Errorf("(error)  command is '%s'", ("nil"))
}

func ErrConnect(s ...string) error {
	return fmt.Errorf("(error)  connect error '%s'", (s[0]))
}

func ErrRead(s ...string) error {
	return fmt.Errorf("(error)  read command error '%s'", (s[0]))
}

func ErrInvalidResp(s ...string) error {
	return fmt.Errorf("(error)  %s", ("invalid response"))
}
