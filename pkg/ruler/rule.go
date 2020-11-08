package ruler

/*
This struct is the main format for rules or conditions in ruler-compatable libraries.
Here's a sample in JSON format:
	{
		"comparator": "eq",
		"path": "person.name",
		"value": "James"
	}

Valid comparators are: eq, neq, lt, lte, gt, gte, contains (regex), ncontains (!regex)

This struct is exported here so that you can include it in your own JSON encoding/decoding,
but go-ruler has a facility to help decode your rules from JSON into its own structs.
*/
type Rule struct {
	Comparator comparator    `json:"comparator"`
	Path       string        `json:"path"`
	Values     []interface{} `json:"values"`
}

type comparator string

const (
	eq  comparator = "eq"
	neq comparator = "neq"
	gt  comparator = "gt"
	gte comparator = "gte"
	lt  comparator = "lt"
	lte comparator = "lte"
	in  comparator = "in"
)

var allComparators = []comparator{
	eq, neq,
	gt, gte,
	lt, lte,
	in,
}

func (c comparator) Valid() bool {
	for _, v := range allComparators {
		if c == v {
			return true
		}
	}

	return false
}

// Implements the stringer interface
func (c comparator) String() string {
	return string(c)
}
