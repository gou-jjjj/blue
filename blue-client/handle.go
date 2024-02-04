package main

import (
	"fmt"
	"strconv"
)

func Exec(v []string) ([]byte, error) {
	header, ok := bsp.HandleMap2[fmt.Sprintf("%s %s", v[0], v[1])]
	if !ok {
		return nil, ErrCommandType(v...)
	}

	switch v[0] {
	case "sys":
		return bsp.NewBspSysReq(header), nil
	case "num":
		if len(v) <= 3 {
			return bsp.NewBspDataReq(header, v[2]), nil
		}
		number, err := str2Number(v[3])
		if err != nil {
			return nil, err
		}
		return bsp.NewBspDataReq(header, v[2], number), nil
	case "str":
		val := make([]byte, 0, 512)
		for _, s := range v[3:] {
			val = append(val, []byte(s)...)
			val = append(val, []byte("\n")...)
		}
	case "list":

	default:

	}

	return bsp.NewBspDataReq(header, v[2]), nil
}

func str2Number(s string) ([]byte, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return common.Uint64ToBytes(i), nil
}
