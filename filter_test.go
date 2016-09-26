package filter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

var filters = []struct {
	file     string
	sqlParts map[string]int
	params   []interface{}
}{
	{
		"examples/simple.json",
		map[string]int{
			"AND":                       2,
			"t.id = :tid":               1,
			"s.name = :sname":           1,
			"u.last_name = :ulast_name": 1,
		},
		[]interface{}{1, "Open", "Doe"},
	},
	{
		"examples/operators.json",
		map[string]int{
			"AND":                        7,
			"p.vat < :pvat":              1,
			"c.id IN(:cid_0,:cid_1,:cid_2,:cid_3)":               1,
			"u.state NOT IN(:ustate_0,:ustate_1,:ustate_2)":     1,
			"s = :s":                     1,
			"u.last_name <> :ulast_name": 1,
			"min_rate > :min_rate":       1,
			"max_rate >= :max_rate":      1,
		},
		[]interface{}{1, "Open", "Doe"},
	},
}

func TestEqual(t *testing.T) {

	filter := NewFilter("t", map[string]string{
		"minRate":       "min_rate",
		"maxRate":       "max_rate",
		"user":          "u",
		"user.lastName": "u.last_name",
		"state":         "s",
		"customer":      "c",
		"project":       "p",
	})

	for _, tt := range filters {

		b, err := ioutil.ReadFile(tt.file)
		assert.Nil(t, err)

		node, err := filter.Parse(string(b))
		assert.Nil(t, err)
		assert.NotNil(t, node)

		visitor := &SQLVisitor{}
		node.Accept(visitor)

		sql := visitor.Sql()
		fmt.Println(sql)
		// golang maps are not ordered ..
		for part, count := range tt.sqlParts {
			if count == 1 {
				assert.Contains(t, sql, part)
			} else {
				assert.Equal(t, count, strings.Count(sql, part))
			}
		}

		//assert.Contains()
	}

}
