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
}

type CompareNode struct {
	Field    string
	Value    interface{}
	Operator uint32
}

func NewCompareNode(field string, value interface{}, operator uint32) (*CompareNode, error) {
	if field == "" {
		return nil, errors.New("Fieldname may not be empty")
	}
	switch operator {
	case compareEq:
	case compareNotEq:
	case compareGt:
	case compareGte:
	case compareLt:
	case compareLte:
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
	Operator  uint8
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

func (e *LogicNode) Accept(visitor Visitor) {
	if e.LeftNode != nil {
		e.LeftNode.Accept(visitor)
	}
	visitor.Visit(e)
	if e.RightNode != nil {
		e.RightNode.Accept(visitor)
	}
}
