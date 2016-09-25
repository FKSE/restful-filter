package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

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
		return nil, err
	}
	return f.parseNode(q)
}

func (f *Filter) parseNode(node map[string]interface{}) (Node, error) {
	var tree Node
	fmt.Println(node)
	for k, v := range node {
		// check if field is a operator
		if operator, ok := compareOperatorMapping[k]; ok {
			fmt.Println(operator)
		} else {
			var node Node
			var err error
			// either implicit and or
			switch v.(type) {
			case string, bool, float64, nil:
				node, err = NewCompareNode(f.replaceAlias(k), v, compareEq)
			case []interface{}:
				node, err = NewCompareNode(f.replaceAlias(k), v, compareIn)
			case map[string]interface{}:
			}
			if err != nil {
				return nil, err
			}
			//fmt.Println(node)
			if tree == nil {
				tree = node
			} else {
				tree = f.insertNode(tree, node)
			}
		}
	}
	return tree, nil
}

func (f *Filter) insertNode(tree Node, node Node) Node {
	if logicNode, ok := tree.(*LogicNode); ok {
		if logicNode.LeftNode == nil {
			logicNode.LeftNode = node
		} else if logicNode.RightNode == nil {
			logicNode.RightNode = node
		} else {
			tree, _ = NewLogicNode(tree, node, logicNode.Operator)
		}
	} else {
		tree, _ = NewLogicNode(tree, node, logicAnd)
	}
	return tree
}

func (f *Filter) replaceAlias(field string) string {
	if repl, ok := f.mapping[field]; ok {
		return repl
	}
	index := strings.LastIndex(field, ".")
	if index == -1 && f.rootAlias != "" {
		field = fmt.Sprintf("%s.%s", f.rootAlias, field)
	} else if f.mapping != nil {
		if prefix, ok := f.mapping[field[0:index]]; ok {
			field = strings.Replace(field, field[0:index], prefix, 1)
		}
	}
	return field
}
