package set

import (
	iter "blue/datastruct"
	mapset "github.com/deckarep/golang-set/v2"
	"strings"
)

type Set struct {
	mapset.Set[string]
	iter.BlueObj
}

func NewSet() Set {
	return Set{
		Set: mapset.NewSet[string](),
		BlueObj: iter.BlueObj{
			Type: iter.Set,
		},
	}
}

func (s *Set) Type() string {
	return s.GetType()
}

func (s *Set) String() string {

	builder := strings.Builder{}
	s.Set.Each(func(s string) bool {
		builder.WriteString(s)
		builder.WriteRune(' ')
		return false
	})

	res := builder.String()[:builder.Len()-1]
	return res
}
