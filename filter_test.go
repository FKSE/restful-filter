package filter

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var filters = []struct {
	file   string
	sql    string
	params []interface{}
}{
	{
		"examples/simple.json",
		"t.id = ? AND s.name = ? AND u.last_name = ?",
		[]interface{}{1, "Open", "Doe"},
	},
}

func TestEqual(t *testing.T) {

	filter := NewFilter("t", map[string]string{
		"user":  "u",
		"state": "s",
	})

	for _, tt := range filters {

		b, err := ioutil.ReadFile(tt.file)
		assert.Nil(t, err)

		node, err := filter.Parse(string(b))
		assert.Nil(t, err)

		visitor := SQLVisitor{}
		node.Accept(visitor)
	}

}
