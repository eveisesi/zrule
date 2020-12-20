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

func Validate(rules [][]*Rule) error {

	for _, andRules := range rules {
		for _, rule := range andRules {
			if err := rule.Validate(); err != nil {
				return err
			}
		}
	}

	return nil

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
func (r *Ruler) Test(o interface{}) (bool, error) {

	for _, andRule := range r.rules {
		var passedCounter int
		for _, rule := range andRule {

			values, err := r.ValuesToEvaluate(rule.Path, 0, reflect.ValueOf(o), make(map[interface{}]bool)))
			if err != nil {
				return false, err
			}

			for _, ruleValue := range rule.Values {

				for _, value := range values {

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

}


// ValuesToEvaluate takes in a path, initial depth, and reflect.Value and return values that can be evaluated based on the path
// The path is a dotted string (a.b.e) and v can be any complex go value, struct, map, array, etc.
// Credit: https://stackoverflow.com/questions/47664320/golang-recursively-reflect-both-type-of-field-and-value#answer-47664689
func (r *Ruler) ValuesToEvaluate(path string, depth int, v reflect.Value, visited map[interface{}]bool) []interface{} {

	results := make([]interface{}, 0)

	parts := strings.Split(path, ".")
	if len(parts) < depth {
		return results
	}

	// Drill down through pointers and interfaces to get a value we can print.
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.Kind() == reflect.Ptr {
			// Check for recursive data
			if visited[v.Interface()] {
				return results
			}
			visited[v.Interface()] = true
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			results = append(results, values(path, depth, v.Index(i), visited)...)
		}
	case reflect.Struct:
		key := parts[depth]
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).Name != key {
				continue
			}
			results = append(results, values(path, depth+1, v.Field(i), visited)...)
		}
	case reflect.Map:
		key := parts[depth]
		for _, e := range v.MapKeys() {
			if e.String() != key {
				continue
			}
			results = append(results, values(path, depth+1, v.MapIndex(e), visited)...)
		}

	default:
		results = append(results, v.Interface())
		return results
	}

	return results
}

func (r *Ruler) EvaluateValue(cp Comparator, ruleValue interface{}, value interface{}) (bool, error) {
	// Grab the type of the value from the Rule and the type of the value that we are comparing to the rule

	switch cp {
	case EQ, IN:
		return reflect.DeepEqual(ruleValue, value), nil
	case NEQ:
		return !reflect.DeepEqual(ruleValue, value), nil
	case GT, GTE, LT, LTE:
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
func (r *Ruler) inequality(op Comparator, actual, expected interface{}) bool {
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

func compareStrings(op Comparator, actual, expected string) bool {
	switch op {
	case GT:
		return actual > expected
	case GTE:
		return actual >= expected
	case LT:
		return actual < expected
	case LTE:
		return actual <= expected
	default:
		return false
	}
}

func compareInts(op Comparator, actual, expected int64) bool {
	switch op {
	case GT:
		return actual > expected
	case GTE:
		return actual >= expected
	case LT:
		return actual < expected
	case LTE:
		return actual <= expected
	default:
		return false
	}
}

func compareUints(op Comparator, actual, expected uint64) bool {
	switch op {
	case GT:
		return actual > expected
	case GTE:
		return actual >= expected
	case LT:
		return actual < expected
	case LTE:
		return actual <= expected
	default:
		return false
	}
}

func compareFloats(op Comparator, actual, expected float64) bool {
	switch op {
	case GT:
		return actual > expected
	case GTE:
		return actual >= expected
	case LT:
		return actual < expected
	case LTE:
		return actual <= expected
	default:
		return false
	}
}
