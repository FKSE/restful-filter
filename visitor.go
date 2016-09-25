package filter

import "fmt"

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
	v.sql += fmt.Sprintf("%s %s ?", node.Field, sqlOperators[node.Operator])
}

func (v *SQLVisitor) Sql() string {
	return v.sql
}
