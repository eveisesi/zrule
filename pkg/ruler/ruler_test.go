package ruler

import (
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
			[]*Rule{
				&Rule{
					"eq",
					"basic.property",
					"foobar",
				},
			},
			map[string]interface{}{
				"basic": map[string]interface{}{
					"property": "foobar",
				},
			},

			"testing basic property equality (string)",
			true,
		},
		{
			[]*Rule{
				&Rule{
					"eq",
					"basic.property",
					12,
				},
			},
			map[string]interface{}{
				"basic": map[string]interface{}{
					"property": 12,
				},
			},
			"testing basic property equality (int)",
			true,
		},
		{
			[]*Rule{
				&Rule{
					"gt",
					"basic.property",
					45,
				},
			},
			map[string]interface{}{
				"basic": map[string]interface{}{
					"property": 100,
				},
			},
			"testing greater than (int)",
			true,
		},
		{
			[]*Rule{
				&Rule{
					"gte",
					"basic.property",
					100,
				},
			},
			map[string]interface{}{
				"basic": map[string]interface{}{
					"property": 100,
				},
			},
			"testing greater than or equal to (int)",
			true,
		},
		{
			[]*Rule{
				&Rule{
					"lt",
					"basic.property",
					45,
				},
			},
			map[string]interface{}{
				"basic": map[string]interface{}{
					"property": 10,
				},
			},
			"testing less than (int)",
			true,
		},
		{
			[]*Rule{
				&Rule{
					"lte",
					"basic.property",
					45,
				},
			},
			map[string]interface{}{
				"basic": map[string]interface{}{
					"property": 45,
				},
			},
			"testing less than or equal to (int)",
			true,
		},
		{
			[]*Rule{
				&Rule{
					"eq",
					"basic.property",
					42,
				},
			},
			map[string]interface{}{
				"basic": []interface{}{
					map[string]interface{}{
						"property": 42,
					},
				},
			},
			"testing equality on an array with truth as first index array (int)",
			true,
		},
		{
			[]*Rule{
				&Rule{
					"eq",
					"basic.property",
					42,
				},
			},
			map[string]interface{}{
				"basic": []interface{}{
					map[string]interface{}{
						"property": 32,
					},
					map[string]interface{}{
						"property": 92,
					},
				},
			},
			"testing equality on an array with truth in second index (int)",
			false,
		},
		{
			[]*Rule{
				&Rule{
					"eq",
					"basic.property",
					92,
				},
			},
			map[string]interface{}{
				"basic": []interface{}{
					map[string]interface{}{
						"property": 32,
					},
					map[string]interface{}{
						"property": 92,
					},
				},
			},
			"testing equality on an array with truth in second index (int)",
			true,
		},
		{
			[]*Rule{
				&Rule{
					"eq",
					"a.b.c",
					1,
				},
			},
			map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{
						"b": []interface{}{
							map[string]interface{}{
								"c": 1,
							},
						},
					},
				},
			},
			"testing equality on deeply nest array",
			true,
		},
		{
			[]*Rule{
				&Rule{
					"gt",
					"a.b.c",
					2,
				},
			},
			map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{
						"b": []interface{}{
							map[string]interface{}{
								"c": 4,
							},
						},
					},
				},
			},
			"testing inequality on deeply nest array",
			true,
		},
	}

	for _, c := range cases {
		r := &Ruler{
			rules: c.rules,
		}

		result, err := r.Test(c.o)
		if err != nil {
			t.Errorf("rule test failed! %s\nrules: %s",
				c.name,
				c.rules,
			)
		}

		if result != c.result {
			t.Errorf("rule test failed! %s\nrules: %s\nExpected %t, Got %t",
				c.name,
				c.rules,
				c.result,
				result,
			)
		}
	}
}
