package filter

import "encoding/json"

var compareOperatorMapping = map[string]uint32{
	"$eq":  compareEq,
	"$ne":  compareNotEq,
	"$gt":  compareGt,
	"$gte": compareGte,
	"$lt":  compareLt,
	"$lte": compareLte,
	"$in":  compareIn,
	"$nin": compareNotIn,
}

var logicOperatorMapping = map[string]uint32{
	"$and": logicAnd,
	"$or":  logicOr,
	"$not": logicNot,
}

type Filter struct {
	rootAlias string
	mapping   map[string]string
}

func NewFilter(rootAlias string, aliasMap map[string]string) *Filter {
	return &Filter{
		rootAlias: rootAlias,
		mapping:   aliasMap,
	}
}

func (f *Filter) Parse(query string) (Node, error) {
	var q map[string]interface{}
	err := json.Unmarshal([]byte(query), &q)
	if err != nil {
		return err
	}
	var tree Node
	for k, v := range q {
		// check if field is a operator
		if operator, ok := compareOperatorMapping[k]; ok {

		} else {
			var node Node
			var err error
			// either implicit and or
			switch v.(type) {
			case string:
			case bool:
			case float64:
			case nil:
				node, err = NewCompareNode(f.replaceAlias(k), v, compareEq)
			case []interface{}:
				node, err = NewCompareNode(f.replaceAlias(k), v, compareIn)
			case map[string]interface{}:

			}
			if err != nil {
				return nil, err
			}
			f.insertNode(tree, node)
		}
	}
	return nil, nil
}

func (f *Filter) insertNode(tree Node, node Node) {
	if tree == nil {
		tree = node
		return
	}
	if logicNode, ok := tree.(LogicNode); ok {
		if logicNode.LeftNode == nil {
			logicNode.LeftNode = node
		} else if logicNode.RightNode == nil {
			logicNode.RightNode = node
		} else {

		}
	} else {

	}

}

func (f *Filter) replaceAlias(field string) string {
	return field
}
