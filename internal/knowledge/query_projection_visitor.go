package knowledge

import (
	"fmt"

	"github.com/clems4ever/go-graphkb/internal/query"
)

type ProjectionVisitor struct {
	ExpressionVisitorBase

	queryGraph *QueryGraph

	Aggregation    bool
	TypeAndIndex   TypeAndIndex
	ExpressionType ExpressionType

	funcInvoc  bool
	etype      ExpressionType
	properties []string
}

func NewProjectionVisitor(queryGraph *QueryGraph) *ProjectionVisitor {
	return &ProjectionVisitor{
		queryGraph: queryGraph,
	}
}

// ParseExpression return whether the expression require aggregation
func (pv *ProjectionVisitor) ParseExpression(q *query.QueryExpression) error {
	err := NewExpressionParser(pv, pv.queryGraph).ParseExpression(q)
	if err != nil {
		return err
	}
	return nil
}

// OnEnterFunctionInvocation called when the EnterFunctionoInvocation is parsed. Name is the name of the function.
func (pv *ProjectionVisitor) OnEnterFunctionInvocation(name string, distinct bool) error {
	if name == "COUNT" {
		pv.Aggregation = true
		pv.funcInvoc = true
	} else {
		return fmt.Errorf("Function %s is not supported", name)
	}
	return nil
}

func (pv *ProjectionVisitor) OnVariable(name string) error {
	typeAndIndex, err := pv.queryGraph.FindVariable(name)
	if err != nil {
		return err
	}

	switch typeAndIndex.Type {
	case NodeType:
		pv.etype = NodeExprType
	case RelationType:
		pv.etype = EdgeExprType
	default:
		pv.etype = PropertyExprType
	}

	pv.TypeAndIndex = typeAndIndex
	return nil
}

func (pv *ProjectionVisitor) OnVariablePropertiesPath(properties []string) error {
	pv.properties = properties
	return nil
}

func (pv *ProjectionVisitor) OnExitPropertyOrLabelsExpression(e query.QueryPropertyOrLabelsExpression) error {
	if len(pv.properties) > 0 || pv.funcInvoc {
		pv.ExpressionType = PropertyExprType
	} else {
		pv.ExpressionType = pv.etype
	}
	pv.properties = nil
	return nil
}
