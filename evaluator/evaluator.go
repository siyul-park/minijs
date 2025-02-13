package evaluator

import (
	"fmt"
	"github.com/siyul-park/miniscript/ast"
	"github.com/siyul-park/miniscript/token"
	"github.com/siyul-park/miniscript/types"
)

func Eval(node ast.Node, env *Environment) (types.Value, error) {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n, env)
	case *ast.PrefixExpression:
		return evalPrefix(n, env)
	case *ast.InfixExpression:
		return evalInfix(n, env)
	case *ast.IntLiteral:
		return types.NewInt(n.Value), nil
	case *ast.FloatLiteral:
		return types.NewFloat(n.Value), nil
	default:
		return nil, fmt.Errorf("unknown node type: %T", n)
	}
}

func evalProgram(p *ast.Program, env *Environment) (types.Value, error) {
	var value types.Value = types.NULL
	for _, stmt := range p.Nodes {
		var err error
		value, err = Eval(stmt, env)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

func evalPrefix(expr *ast.PrefixExpression, env *Environment) (types.Value, error) {
	right, err := Eval(expr.Right, env)
	if err != nil {
		return nil, err
	}

	switch expr.Token.Type {
	case token.PLUS:
		return right, nil
	case token.MINUS:
		switch v := right.(type) {
		case *types.Int:
			return types.NewInt(-v.Value), nil
		case *types.Float:
			return types.NewFloat(-v.Value), nil
		default:
			return nil, fmt.Errorf("unsupported type for negation: %T", v)
		}
	default:
		return nil, fmt.Errorf("unknown operator: %s", expr.Token.Type)
	}
}

func evalInfix(expr *ast.InfixExpression, env *Environment) (types.Value, error) {
	left, err := Eval(expr.Left, env)
	if err != nil {
		return nil, err
	}
	right, err := Eval(expr.Right, env)
	if err != nil {
		return nil, err
	}

	switch expr.Token.Type {
	case token.PLUS:
		switch left := left.(type) {
		case *types.Int:
			switch right := right.(type) {
			case *types.Int:
				return types.NewInt(left.Value + right.Value), nil
			case *types.Float:
				return types.NewFloat(float64(left.Value) + right.Value), nil
			default:
				return nil, fmt.Errorf("unsupported type for addition: %T", right)
			}
		case *types.Float:
			switch right := right.(type) {
			case *types.Int:
				return types.NewFloat(left.Value + float64(right.Value)), nil
			case *types.Float:
				return types.NewFloat(left.Value + right.Value), nil
			default:
				return nil, fmt.Errorf("unsupported type for addition: %T", right)
			}
		default:
			return nil, fmt.Errorf("unsupported type for addition: %T", left)
		}
	case token.MINUS:
		switch left := left.(type) {
		case *types.Int:
			switch right := right.(type) {
			case *types.Int:
				return types.NewInt(left.Value - right.Value), nil
			case *types.Float:
				return types.NewFloat(float64(left.Value) - right.Value), nil
			default:
				return nil, fmt.Errorf("unsupported type for subtraction: %T", right)
			}
		case *types.Float:
			switch right := right.(type) {
			case *types.Int:
				return types.NewFloat(left.Value - float64(right.Value)), nil
			case *types.Float:
				return types.NewFloat(left.Value - right.Value), nil
			default:
				return nil, fmt.Errorf("unsupported type for subtraction: %T", right)
			}
		default:
			return nil, fmt.Errorf("unsupported type for subtraction: %T", left)
		}
	case token.MULTIPLY:
		switch left := left.(type) {
		case *types.Int:
			switch right := right.(type) {
			case *types.Int:
				return types.NewInt(left.Value * right.Value), nil
			case *types.Float:
				return types.NewFloat(float64(left.Value) * right.Value), nil
			default:
				return nil, fmt.Errorf("unsupported type for multiplication: %T", right)
			}
		case *types.Float:
			switch right := right.(type) {
			case *types.Int:
				return types.NewFloat(left.Value * float64(right.Value)), nil
			case *types.Float:
				return types.NewFloat(left.Value * right.Value), nil
			default:
				return nil, fmt.Errorf("unsupported type for multiplication: %T", right)
			}
		default:
			return nil, fmt.Errorf("unsupported type for multiplication: %T", left)
		}
	case token.DIVIDE:
		switch left := left.(type) {
		case *types.Int:
			switch right := right.(type) {
			case *types.Int:
				if right.Value == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				return types.NewInt(left.Value / right.Value), nil
			case *types.Float:
				if right.Value == 0.0 {
					return nil, fmt.Errorf("division by zero")
				}
				return types.NewFloat(float64(left.Value) / right.Value), nil
			default:
				return nil, fmt.Errorf("unsupported type for division: %T", right)
			}
		case *types.Float:
			switch right := right.(type) {
			case *types.Int:
				if right.Value == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				return types.NewFloat(left.Value / float64(right.Value)), nil
			case *types.Float:
				if right.Value == 0.0 {
					return nil, fmt.Errorf("division by zero")
				}
				return types.NewFloat(left.Value / right.Value), nil
			default:
				return nil, fmt.Errorf("unsupported type for division: %T", right)
			}
		default:
			return nil, fmt.Errorf("unsupported type for division: %T", left)
		}
	default:
		return nil, fmt.Errorf("unknown operator: %s", expr.Token.Type)
	}
}
