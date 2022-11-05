package eval

import (
	"github.com/MichaelDiBernardo/monkey/ast"
	"github.com/MichaelDiBernardo/monkey/object"
	"github.com/MichaelDiBernardo/monkey/token"
)

func Eval(root ast.Node) object.Object {
	switch node := root.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Value)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		if node.Value {
			return object.TRUE_OBJ
		}
		return object.FALSE_OBJ
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

func evalPrefixExpression(node *ast.PrefixExpression) object.Object {
	rhsval := Eval(node.RHS)
	switch node.OperatorToken.Type {
	case token.BANG:
		return evalBangOperatorExpression(rhsval)
	default:
		return object.NULL_OBJ
	}
}

func evalBangOperatorExpression(rhsval object.Object) object.Object {
	switch rhsval {
	case object.TRUE_OBJ:
		return object.FALSE_OBJ
	case object.FALSE_OBJ:
		return object.TRUE_OBJ
	case object.NULL_OBJ:
		return object.TRUE_OBJ
	default:
		return object.FALSE_OBJ
	}
}
