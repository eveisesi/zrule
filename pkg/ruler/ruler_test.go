package ruler_test

import (
	"encoding/json"
	"testing"

	"github.com/eveisesi/zrule"
	"github.com/eveisesi/zrule/pkg/ruler"
)

func TestRules(t *testing.T) {

	cases := []struct {
		rules  ruler.Rules
		o      interface{}
		name   string
		result bool
	}{
		{
			[][]*ruler.Rule{
				[]*ruler.Rule{
					&ruler.Rule{
						"eq",
						"Meta.Solo",
						[]interface{}{true},
					},
				},
			},
			zrule.Killmail{
				Meta: &zrule.Meta{
					Solo: true,
				},
			},
			"testing boolean on nested struct",
			true,
		},
		{
			[][]*ruler.Rule{
				[]*ruler.Rule{
					&ruler.Rule{
						"eq",
						"Victim.ShipTypeID",
						[]interface{}{670},
					},
					&ruler.Rule{
						"gt",
						"Meta.TotalValue",
						[]interface{}{10000},
					},
				},
			},
			zrule.Killmail{
				Victim: &zrule.KillmailVictim{
					ShipTypeID: 670,
				},
				Meta: &zrule.Meta{
					TotalValue: 10000,
				},
			},
			"testing and rule, comparing victim ship and killmail total value, should return false",
			false,
		},
		{
			[][]*ruler.Rule{
				[]*ruler.Rule{
					&ruler.Rule{
						"eq",
						"Victim.ShipTypeID",
						[]interface{}{670},
					},
					&ruler.Rule{
						"gt",
						"Meta.TotalValue",
						[]interface{}{10000},
					},
				},
			},
			zrule.Killmail{
				Victim: &zrule.KillmailVictim{
					ShipTypeID: 670,
				},
				Meta: &zrule.Meta{
					TotalValue: 15000,
				},
			},
			"testing and rule, comparing victim ship and killmail total value, should return true",
			true,
		},
	}

	for _, c := range cases {
		r := ruler.NewRuler()
		r.SetRules(c.rules)

		result := r.Test(c.o)
		if result != c.result {
			values, _ := json.Marshal(c.o)
			t.Errorf("Test Failed:\nName: %s\nRules: %s\nValues: %s\nExpected %t, Got %t",
				c.name,
				c.rules,
				string(values),
				c.result,
				result,
			)
		}
	}
}
