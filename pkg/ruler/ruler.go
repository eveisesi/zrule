package ruler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type Rules [][]*Rule

func (r Rules) String() string {

	data, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(data)
}

type Ruler struct {
	rules  Rules
	config struct {
		ErrOnMissingKey bool
	}
}

func NewRuler() *Ruler {
	return &Ruler{}
}

func OptionErrOnMissingKey(r *Ruler) {
	r.config.ErrOnMissingKey = true
}

// SetRules takes in a slice of rules and set the on the ruler
func (r *Ruler) SetRules(s [][]*Rule) {
	r.rules = s
}

// SetRulesWithJSON takes in a slice of rule, unmarshals them onto a slice of rule, panic if unmarshal errors
func (r *Ruler) SetRulesWithJSON(d []byte) {

	s := make([][]*Rule, 0)
	err := json.NewDecoder(bytes.NewReader(d)).Decode(&s)
	if err != nil {
		panic(err)
	}

	r.rules = s

}

// Test takes in an interface. Underlying type should be a
// map[string]interface or a map[string][]interface{}
// The root element should be the equivalent of a JSON Object
func (r *Ruler) Test(o map[string]interface{}) (bool, error) {

	for _, andRule := range r.rules {
		var passedCounter int
		for _, rule := range andRule {

			values, err := r.ValuesToEvaluate(rule.Path, 0, o)
			if err != nil {
				return false, err
			}

			for _, ruleValue := range rule.Values {

				for _, value := range values {
					fmt.Println(ruleValue, value)

					passed, err := r.EvaluateValue(rule.Comparator, ruleValue, value)
					if err != nil {
						return false, err
					}
					if passed {
						passedCounter++
						goto NextRule
					}

				}

			}
		NextRule:
		}

		if passedCounter == len(andRule) {
			return true, nil
		}

	}

	return false, nil

	// for _, rule := range r.rules {

	// 	// At the top we should get a boolean response out of the interface
	// 	values, err := r.ValuesToEvaluate(rule, 0, o)
	// 	if err != nil {
	// 		return false, err
	// 	}
	// 	var passed bool
	// 	for _, ruleValue := range rule.Values {

	// 		for _, value := range values {
	// 			passed, err = r.EvaluateValue(rule.Comparator, ruleValue, value)
	// 			if err != nil {
	// 				return false, err
	// 			}
	// 			if passed {
	// 				goto AllValuesBreak
	// 			}
	// 		}

	// 	}
	// AllValuesBreak:
	// 	if passed {
	// 		return true, nil
	// 	}

	// }

}

var ErrMissingKeyFmt = "no value found for key %s"

func (r *Ruler) ValuesToEvaluate(path string, depth int, o interface{}) ([]interface{}, error) {

	parts := strings.Split(path, ".")
	if len(parts) <= depth {
		return nil, fmt.Errorf("depth greater than length of key")
	}

	key := parts[depth]
	values := make([]interface{}, 0)

	switch a := o.(type) {
	case map[string]interface{}:
		if _, ok := a[key]; !ok {
			if r.config.ErrOnMissingKey {
				return values, fmt.Errorf(ErrMissingKeyFmt, path)
			}
			return values, nil
		}

		valAtLoc := a[key]
		typeAtLoc := reflect.TypeOf(valAtLoc)

		if typeAtLoc.Kind() == reflect.Slice || typeAtLoc.Kind() == reflect.Map {
			subValues, err := r.ValuesToEvaluate(path, depth+1, valAtLoc)
			if err != nil {
				return values, err
			}

			values = append(values, subValues...)
			return values, nil
		} else if typeAtLoc.Kind() == reflect.Struct {
			return values, fmt.Errorf("Wrapping type must interface or slice interface, Got Struct")
		}

		values = append(values, valAtLoc)
		return values, nil

	case []interface{}:
		for _, b := range a {
			subValues, err := r.ValuesToEvaluate(path, depth, b)
			if err != nil {
				return values, err
			}

			values = append(values, subValues...)
		}
	default:
		return values, fmt.Errorf("unsupported type %T passed in", o)
	}

	return values, nil
}

func (r *Ruler) EvaluateValue(cp comparator, ruleValue interface{}, value interface{}) (bool, error) {
	// Grab the type of the value from the Rule and the type of the value that we are comparing to the rule

	switch cp {
	case eq, in:
		return reflect.DeepEqual(ruleValue, value), nil
	case neq:
		return !reflect.DeepEqual(ruleValue, value), nil
	case gt, gte, lt, lte:
		return r.inequality(cp, value, ruleValue), nil
	default:

		return false, fmt.Errorf("unsupported comparator")
	}

}

// func (r *Ruler) handleIn(expected, actual interface{}) (bool, error) {

// 	return true, nil

// }

// runs equality comparison
// separated in a different function because
// we need to do another type assertion here
// and some other acrobatics
func (r *Ruler) inequality(op comparator, actual, expected interface{}) bool {
	// need some variables for these deals
	var cmpStr [2]string
	var isStr [2]bool
	var cmpUint [2]uint64
	var isUint [2]bool
	var cmpInt [2]int64
	var isInt [2]bool
	var cmpFloat [2]float64
	var isFloat [2]bool

	for idx, i := range []interface{}{actual, expected} {
		switch t := i.(type) {
		case uint8:
			cmpUint[idx] = uint64(t)
			isUint[idx] = true
		case uint16:
			cmpUint[idx] = uint64(t)
			isUint[idx] = true
		case uint32:
			cmpUint[idx] = uint64(t)
			isUint[idx] = true
		case uint64:
			cmpUint[idx] = t
			isUint[idx] = true
		case uint:
			cmpUint[idx] = uint64(t)
			isUint[idx] = true
		case int8:
			cmpInt[idx] = int64(t)
			isInt[idx] = true
		case int16:
			cmpInt[idx] = int64(t)
			isInt[idx] = true
		case int32:
			cmpInt[idx] = int64(t)
			isInt[idx] = true
		case int64:
			cmpInt[idx] = t
			isInt[idx] = true
		case int:
			cmpInt[idx] = int64(t)
			isInt[idx] = true
		case float32:
			cmpFloat[idx] = float64(t)
			isFloat[idx] = true
		case float64:
			cmpFloat[idx] = t
			isFloat[idx] = true
		case string:
			cmpStr[idx] = t
			isStr[idx] = true
		default:
			return false
		}
	}

	if isStr[0] && isStr[1] {
		return compareStrings(op, cmpStr[0], cmpStr[1])
	}

	if isInt[0] && isInt[1] {
		return compareInts(op, cmpInt[0], cmpInt[1])
	}

	if isUint[0] && isUint[1] {
		return compareUints(op, cmpUint[0], cmpUint[1])
	}

	if isFloat[0] && isFloat[1] {
		return compareFloats(op, cmpFloat[0], cmpFloat[1])
	}

	return false
}

func compareStrings(op comparator, actual, expected string) bool {
	switch op {
	case gt:
		return actual > expected
	case gte:
		return actual >= expected
	case lt:
		return actual < expected
	case lte:
		return actual <= expected
	default:
		return false
	}
}

func compareInts(op comparator, actual, expected int64) bool {
	switch op {
	case gt:
		return actual > expected
	case gte:
		return actual >= expected
	case lt:
		return actual < expected
	case lte:
		return actual <= expected
	default:
		return false
	}
}

func compareUints(op comparator, actual, expected uint64) bool {
	switch op {
	case gt:
		return actual > expected
	case gte:
		return actual >= expected
	case lt:
		return actual < expected
	case lte:
		return actual <= expected
	default:
		return false
	}
}

func compareFloats(op comparator, actual, expected float64) bool {
	switch op {
	case gt:
		return actual > expected
	case gte:
		return actual >= expected
	case lt:
		return actual < expected
	case lte:
		return actual <= expected
	default:
		return false
	}
}
