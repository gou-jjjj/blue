package set

import (
	iter "blue/datastruct"
	mapset "github.com/deckarep/golang-set/v2"
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
