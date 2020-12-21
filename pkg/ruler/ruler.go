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
	rules Rules
}

func NewRuler() *Ruler {
	return &Ruler{}
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
func (r *Ruler) Test(o interface{}) bool {

	for _, andRule := range r.rules {
		var passedCounter int
		for _, rule := range andRule {

			values := r.ValuesToEvaluate(rule.Path, 0, reflect.ValueOf(o), make(map[interface{}]bool))
			for _, ruleValue := range rule.Values {

				for _, value := range values {
					passed := r.EvaluateValue(rule.Comparator, value, ruleValue)

					if passed {
						passedCounter++
						goto NextRule
					}

				}

			}

		NextRule:
		}

		if passedCounter == len(andRule) {
			return true
		}

	}

	return false

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
			results = append(results, r.ValuesToEvaluate(path, depth, v.Index(i), visited)...)
		}
	case reflect.Struct:
		key := parts[depth]
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).Name != key {
				continue
			}
			results = append(results, r.ValuesToEvaluate(path, depth+1, v.Field(i), visited)...)
		}
	case reflect.Map:
		key := parts[depth]
		for _, e := range v.MapKeys() {
			if e.String() != key {
				continue
			}
			results = append(results, r.ValuesToEvaluate(path, depth+1, v.MapIndex(e), visited)...)
		}

	default:
		results = append(results, v.Interface())
		return results
	}

	return results
}

// Evaluate Value compares the value received against the expected value in
// rule using the rules registered comparator.
func (r *Ruler) EvaluateValue(op Comparator, actual, expected interface{}) bool {
	var cmpStr [2]string
	var isStr [2]bool
	var cmpFloat [2]float64
	var isFloat [2]bool

	for idx, i := range []interface{}{actual, expected} {
		switch t := i.(type) {
		case uint8:
			cmpFloat[idx] = float64(t)
			isFloat[idx] = true
		case uint16:
			cmpFloat[idx] = float64(t)
			isFloat[idx] = true
		case uint32:
			cmpFloat[idx] = float64(t)
			isFloat[idx] = true
		case uint64:
			cmpFloat[idx] = float64(t)
			isFloat[idx] = true
		case uint:
			cmpFloat[idx] = float64(t)
			isFloat[idx] = true
		case int8:
			cmpFloat[idx] = float64(t)
			isFloat[idx] = true
		case int16:
			cmpFloat[idx] = float64(t)
			isFloat[idx] = true
		case int32:
			cmpFloat[idx] = float64(t)
			isFloat[idx] = true
		case int64:
			cmpFloat[idx] = float64(t)
			isFloat[idx] = true
		case int:
			cmpFloat[idx] = float64(t)
			isFloat[idx] = true
		case float32:
			cmpFloat[idx] = float64(t)
			isFloat[idx] = true
		case float64:
			cmpFloat[idx] = t
			isFloat[idx] = true
		case bool:
			cmpFloat[idx] = float64(0)
			if t {
				cmpFloat[idx] = float64(1)
			}
			isFloat[idx] = true
		case string:
			cmpStr[idx] = t
			isStr[idx] = true
		default:
			panic(fmt.Sprintf("EvaluateValue: Unable to coerve %v (%T) to a float64 or string for comparison", t, t))
		}
	}

	if isStr[0] && isStr[1] {
		return compareStrings(op, cmpStr[0], cmpStr[1])
	}

	if isFloat[0] && isFloat[1] {
		return compareFloats(op, cmpFloat[0], cmpFloat[1])
	}

	return false
}

func compareStrings(op Comparator, actual, expected string) bool {
	switch op {
	case EQ:
		return actual == expected
	case NEQ:
		return actual != expected
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
	case EQ:
		return actual == expected
	case NEQ:
		return actual != expected
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
