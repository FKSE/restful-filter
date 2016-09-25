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
			"AND":             2,
			"t.id = ?":        1,
			"s.name = ?":      1,
			"u.last_name = ?": 1,
		},
		[]interface{}{1, "Open", "Doe"},
	},
}

func TestEqual(t *testing.T) {

	filter := NewFilter("t", map[string]string{
		"user":  "u",
		"user.lastName": "u.last_name",
		"state": "s",
	})

	for _, tt := range filters {

		b, err := ioutil.ReadFile(tt.file)
		assert.Nil(t, err)

		node, err := filter.Parse(string(b))
		assert.Nil(t, err)
		assert.NotNil(t, node)

		visitor := &SQLVisitor{}
		node.Accept(visitor)
		fmt.Println()

		sql := visitor.Sql()
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
