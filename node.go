package filter

import (
	"errors"
	"fmt"
)

const (
	compareEq uint32 = iota
	compareNotEq
	compareGt
	compareGte
	compareLt
	compareLte
	compareIn
	compareNotIn
	logicAnd
	logicOr
	logicNot
)

type Node interface {
	Accept(visitor Visitor)
	String() string
}

type CompareNode struct {
	Field    string
	Value    interface{}
	Operator uint32
}

func (n *CompareNode) Accept(visitor Visitor) {
	visitor.Visit(n)
}

func (n *CompareNode) String() string {
	return fmt.Sprintf("{field:%s,value:%v,operator:%d}", n.Field, n.Value, n.Operator)
}

func NewCompareNode(field string, value interface{}, operator uint32) (*CompareNode, error) {
	if field == "" {
		return nil, errors.New("Fieldname may not be empty")
	}
	switch operator {
	case compareEq, compareNotEq, compareGt, compareGte, compareLt, compareLte:
		return &CompareNode{
			Field:    field,
			Value:    value,
			Operator: operator}, nil
	}
	return nil, fmt.Errorf("Operator %s is not known", operator)
}

type LogicNode struct {
	LeftNode  Node
	RightNode Node
	Operator  uint32
}

func NewLogicNode(left, right Node, operator uint32) (*LogicNode, error) {
	if operator != logicAnd && operator != logicOr && operator != logicNot {
		return nil, errors.New("Unkown operator used")
	}
	return &LogicNode{
		LeftNode:  left,
		RightNode: right,
		Operator:  operator,
	}, nil
}

func (n *LogicNode) String() string {
	node := "{\n\tleft:"
	if n.LeftNode != nil {
		node += n.LeftNode.String()
	} else {
		node += "nil"
	}
	node += fmt.Sprintf("\n\top:%d\n\tright:", n.Operator)
	if n.RightNode != nil {
		node += n.RightNode.String()
	} else {
		node += "nil"
	}
	return node + "\n}"
}

func (n *LogicNode) Accept(visitor Visitor) {
	if n.LeftNode != nil {
		n.LeftNode.Accept(visitor)
	}
	visitor.Visit(n)
	if n.RightNode != nil {
		n.RightNode.Accept(visitor)
	}
}
