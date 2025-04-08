package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"strings"
)

// FunctionEvaluator handles parsing and evaluating mathematical functions
type FunctionEvaluator struct {
	expression string
}

// NewFunctionEvaluator creates a new function evaluator
func NewFunctionEvaluator(function string) *FunctionEvaluator {
	return &FunctionEvaluator{
		expression: function,
	}
}

// Evaluate calculates the result of the function for given x, y values
func (f *FunctionEvaluator) Evaluate(x, y float64) (float64, error) {
	// Replace variables with values
	expr := strings.ReplaceAll(f.expression, "x", fmt.Sprintf("%f", x))
	expr = strings.ReplaceAll(expr, "y", fmt.Sprintf("%f", y))
	
	// Parse the expression
	tree, err := parser.ParseExpr(expr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse expression: %v", err)
	}
	
	// Evaluate the expression
	result, err := f.evaluateExpression(tree)
	if err != nil {
		return 0, err
	}
	
	return result, nil
}

// evaluateExpression recursively evaluates an AST expression
func (f *FunctionEvaluator) evaluateExpression(expr ast.Expr) (float64, error) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		// Handle numeric literals
		if e.Kind == token.FLOAT || e.Kind == token.INT {
			var value float64
			_, err := fmt.Sscanf(e.Value, "%f", &value)
			if err != nil {
				return 0, fmt.Errorf("failed to parse number: %v", err)
			}
			return value, nil
		}
		return 0, fmt.Errorf("unsupported literal type: %v", e.Kind)
	
	case *ast.BinaryExpr:
		// Handle binary operations (+, -, *, /, etc.)
		left, err := f.evaluateExpression(e.X)
		if err != nil {
			return 0, err
		}
		
		right, err := f.evaluateExpression(e.Y)
		if err != nil {
			return 0, err
		}
		
		switch e.Op {
		case token.ADD:
			return left + right, nil
		case token.SUB:
			return left - right, nil
		case token.MUL:
			return left * right, nil
		case token.QUO:
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return left / right, nil
		default:
			return 0, fmt.Errorf("unsupported binary operator: %v", e.Op)
		}
	
	case *ast.ParenExpr:
		// Handle parentheses
		return f.evaluateExpression(e.X)
	
	case *ast.CallExpr:
		// Handle function calls (sin, cos, etc.)
		fun, ok := e.Fun.(*ast.Ident)
		if !ok {
			return 0, fmt.Errorf("unsupported function call")
		}
		
		if len(e.Args) != 1 {
			return 0, fmt.Errorf("function %s requires exactly one argument", fun.Name)
		}
		
		arg, err := f.evaluateExpression(e.Args[0])
		if err != nil {
			return 0, err
		}
		
		switch fun.Name {
		case "sin":
			return math.Sin(arg), nil
		case "cos":
			return math.Cos(arg), nil
		case "tan":
			return math.Tan(arg), nil
		case "sqrt":
			if arg < 0 {
				return 0, fmt.Errorf("square root of negative number")
			}
			return math.Sqrt(arg), nil
		case "exp":
			return math.Exp(arg), nil
		case "log":
			if arg <= 0 {
				return 0, fmt.Errorf("logarithm of non-positive number")
			}
			return math.Log(arg), nil
		default:
			return 0, fmt.Errorf("unsupported function: %s", fun.Name)
		}
	
	case *ast.UnaryExpr:
		// Handle unary operations (-x, +x)
		operand, err := f.evaluateExpression(e.X)
		if err != nil {
			return 0, err
		}
		
		switch e.Op {
		case token.SUB:
			return -operand, nil
		case token.ADD:
			return operand, nil
		default:
			return 0, fmt.Errorf("unsupported unary operator: %v", e.Op)
		}
	
	default:
		return 0, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

// GeneratePointsFromFunction creates a set of 3D points based on the provided function
// and adds them to the given Space3D instance
func GeneratePointsFromFunction(space *Space3D, function string, xMin, xMax, yMin, yMax, step float64) error {
	// Prepare function for evaluation
	eval := NewFunctionEvaluator(function)
	
	// Generate grid of points
	for x := xMin; x <= xMax; x += step {
		for y := yMin; y <= yMax; y += step {
			// Evaluate function to get z value
			z, err := eval.Evaluate(x, y)
			if err != nil {
				// Skip points where evaluation fails
				continue
			}
			
			// Add point to the space
			space.AddPoint(NewPoint3D(x, z, y)) // Note: Using z as the y-coordinate for better visualization
		}
	}
	
	return nil
}