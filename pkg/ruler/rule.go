package ruler

import "fmt"

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
	Comparator Comparator    `bson:"comparator" json:"comparator"`
	Path       string        `bson:"path" json:"path"`
	Values     []interface{} `bson:"values" json:"values"`
}

func (r Rule) Validate() error {

	if !r.Comparator.Valid() {
		return fmt.Errorf("invalid comparator %s specified", r.Comparator)
	}
	if len(r.Path) == 0 {
		return fmt.Errorf("empty path specified, please specific valid path")
	}
	if len(r.Values) == 0 {
		return fmt.Errorf("no rule values specified. Please specific atleast one value for the rule to match against")
	}

	if len(r.Values) > 1 && r.Comparator != IN {
		return fmt.Errorf("invalid values specified for provided comparator. Comparator must be in when values greater than 1")
	}

	return nil
}

type Comparator string

const (
	EQ  Comparator = "eq"
	NEQ Comparator = "neq"
	GT  Comparator = "gt"
	GTE Comparator = "gte"
	LT  Comparator = "lt"
	LTE Comparator = "lte"
	IN  Comparator = "in"
)

var AllComparators = []Comparator{
	EQ, NEQ,
	GT, GTE,
	LT, LTE,
	IN,
}

func (c Comparator) Valid() bool {
	for _, v := range AllComparators {
		if c == v {
			return true
		}
	}

	return false
}

// Implements the stringer interface
func (c Comparator) String() string {
	return string(c)
}
