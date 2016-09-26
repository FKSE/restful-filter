package filter

import (
	"fmt"
	"strings"
)

type Visitor interface {
	Visit(node Node)
}

var sqlOperators = map[uint32]string{
	compareEq:    "=",
	compareNotEq: "<>",
	compareGt:    ">",
	compareGte:   ">=",
	compareLt:    "<",
	compareLte:   "<=",
	compareIn:    "IN",
	compareNotIn: "NOT IN",
	logicAnd:     "AND",
	logicOr:      "OR",
	logicNot:     "NOT",
}

type SQLVisitor struct {
	sql string
}

func (v *SQLVisitor) Visit(node Node) {
	if logicNode, ok := node.(*LogicNode); ok {
		v.visitLogic(logicNode)
	} else {
		v.visitCompare(node.(*CompareNode))
	}
}

func (v *SQLVisitor) visitLogic(node *LogicNode) {
	v.sql += " " + sqlOperators[node.Operator] + " "
}

func (v *SQLVisitor) visitCompare(node *CompareNode) {
	param := strings.Replace(node.Field, ".", "", -1)
	if node.Operator == compareIn || node.Operator == compareNotIn {

		var params []string
		for index := range node.Value.([]interface{}) {
			params = append(params, fmt.Sprintf(":%s_%d", param, index))
		}
		v.sql += fmt.Sprintf("%s %s(%s)", node.Field, sqlOperators[node.Operator], strings.Join(params, ","))
	} else {
		v.sql += fmt.Sprintf("%s %s :%s", node.Field, sqlOperators[node.Operator], param)
	}
}

func (v *SQLVisitor) Sql() string {
	return v.sql
}
