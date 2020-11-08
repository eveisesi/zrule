package ruler

import (
	"encoding/json"
	"testing"
)

func TestRules(t *testing.T) {

	cases := []struct {
		rules  Rules
		o      map[string]interface{}
		name   string
		result bool
	}{

		{
			[][]*Rule{
				[]*Rule{
					&Rule{
						"in",
						"a.b.c",
						[]interface{}{2, 4, 5, 9},
					},
					&Rule{
						"in",
						"a.b.c",
						[]interface{}{1, 3, 6},
					},
				},
				[]*Rule{
					&Rule{
						"in",
						"a.b.c",
						[]interface{}{9, 56, 46},
					},
					&Rule{
						"in",
						"a.f.g",
						[]interface{}{42, 35, 78},
					},
				},
			},
			map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{
						"b": []interface{}{
							map[string]interface{}{
								"c": 9,
							},
							map[string]interface{}{
								"d": 15,
							},
						},
					},
					map[string]interface{}{
						"f": []interface{}{
							map[string]interface{}{
								"f": 9,
							},
							map[string]interface{}{
								"g": 42,
							},
						},
					},
				},
			},
			"testing of or, first should fail, second should pass",
			true,
		},
	}

	for _, c := range cases {
		r := &Ruler{
			rules: c.rules,
		}

		result, err := r.Test(c.o)
		if err != nil {
			t.Errorf("Test Failed with Error:\nName: %s\nError: %s", c.name, err)
		}

		if result != c.result {
			values, _ := json.Marshal(c.o)
			t.Errorf("Test Failed without Error:\nName: %s\nRules: %s\nValues: %s\nExpected %t, Got %t",
				c.name,
				c.rules,
				string(values),
				c.result,
				result,
			)
		}
	}
}
