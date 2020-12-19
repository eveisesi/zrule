package search

type Key string

const (
	KeyRegions        Key = "regions"
	KeyConstellations Key = "constellations"
	KeySystems        Key = "systems"
	KeyItems          Key = "items"
	KeyItemGroups     Key = "itemGroups"
)

var AllKeys = []Key{
	KeyRegions, KeyConstellations, KeySystems, KeyItemGroups, KeyItems,
}

func (r Key) String() string {
	return string(r)
}

func (r Key) Valid() bool {
	for _, v := range AllKeys {
		if v == r {
			return true
		}
	}
	return false
}
