package eval

import (
	"github.com/MichaelDiBernardo/monkey/ast"
	"github.com/MichaelDiBernardo/monkey/object"
)

func Eval(root ast.Node) object.Object {
	switch node := root.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Value)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}
	return nil
}

func evalStatements(statements []ast.Statement) object.Object {
	var val object.Object = &object.Null{}

	for _, s := range statements {
		val = Eval(s)
	}

	return val
}
