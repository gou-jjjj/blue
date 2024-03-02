package string

import iter "blue/datastruct"

type String_ struct {
	iter.BlueObj
	s string
}

func NewString(str ...string) *String_ {
	s := &String_{
		BlueObj: iter.BlueObj{
			Type: iter.String,
		},
	}

	if len(str) > 0 {
		s.Set(str[0])
	}

	return s
}

func (S *String_) Set(s string) {
	S.s = s
}

func (S *String_) Get() string {
	return S.s
}

func (S *String_) Append(s string) {
	S.s += s
}

func (S *String_) Len() int {
	return len(S.s)
}

func (S *String_) Value() string {
	return S.s
}
