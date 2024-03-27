package set

import mapset "github.com/deckarep/golang-set/v2"

type Set mapset.Set[string]

func NewSet() Set {
	return mapset.NewSet[string]()
}
