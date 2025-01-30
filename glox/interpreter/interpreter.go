package interpreter

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/madraceee/interpreters/glox/environment"
	"github.com/madraceee/interpreters/glox/parser"
	"github.com/madraceee/interpreters/glox/token"
	"github.com/madraceee/interpreters/glox/utils"
)

type Interpreter struct {
	Environment *environment.Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		Environment: environment.NewEnvironment(nil),
	}
}

type RuntimeError struct {
	Token token.Token
	Err   error
}

func (e *RuntimeError) Error() string {
	utils.HadRunTimeError = true
	lineStr := strconv.Itoa(e.Token.Line)
	return "Runtime Error: " + e.Err.Error() + "\n[Line " + lineStr + "]\n"
}

func (i *Interpreter) Interpret(stmts []parser.Stmt) {
	for _, stmt := range stmts {
		_, err := i.execute(stmt)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func (i *Interpreter) execute(stmt parser.Stmt) (interface{}, error) {
	return stmt.Visit(i)
}

func (i *Interpreter) VisitBlockStmt(b *parser.Block) (interface{}, error) {
	return nil, i.executeBlock(b.Statements, environment.NewEnvironment(i.Environment))
}

func (i *Interpreter) VisitVarStmt(v *parser.Var) (interface{}, error) {
	var value interface{}
	var err error
	if v.Initializer != nil {
		value, err = i.evaluate(v.Initializer)
		if err != nil {
			return nil, err
		}
	}

	i.Environment.Define(v.Name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) VisitExpressionStmt(e *parser.Expression) (interface{}, error) {
	_, err := i.evaluate(e.Expression)
	return nil, err
}

func (i *Interpreter) VisitPrintStmt(p *parser.Print) (interface{}, error) {
	val, err := i.evaluate(p.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Println(token.GetStringValue(val.(token.Object)))
	return nil, nil
}

func (i *Interpreter) VisitLiteralExpr(l *parser.Literal) (interface{}, error) {
	return l.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(g *parser.Grouping) (interface{}, error) {
	return i.evaluate(g)
}

func (i *Interpreter) VisitUnaryExpr(u *parser.Unary) (interface{}, error) {
	right, err := i.evaluate(u.Right)
	if err != nil {
		return nil, err
	}

	switch u.Operator.TokenType {
	case token.MINUS:
		err := i.checkNumberOperand(u.Operator, u.Right)
		if err != nil {
			return nil, err
		}
		return token.Object{
			Value_float: -right.(token.Object).Value_float,
			ObjType:     token.NUMBER_TYPE,
		}, nil
	case token.BANG:
		return i.isTruthy(u.Right)
	}

	return nil, &RuntimeError{
		Token: u.Operator,
		Err:   errors.New("undefined unary token"),
	}
}

func (i *Interpreter) VisitAssignExpr(a *parser.Assign) (interface{}, error) {
	val, err := i.evaluate(a.Value)
	if err != nil {
		return nil, err
	}

	err = i.Environment.Assign(a.Name, val)
	return val, err
}

func (i *Interpreter) VisitLogicalExpr(l *parser.Logical) (interface{}, error) {
	left, err := i.evaluate(l.Left)
	if err != nil {
		return nil, err
	}

	val, err := i.isObjectTruthy(left.(token.Object))
	if err != nil {
		return nil, err
	}
	if l.Operator.TokenType == token.OR {
		if val.Value_bool {
			return left, nil
		}
	} else {
		if !val.Value_bool {
			return left, nil
		}
	}

	return i.evaluate(l.Right)
}

func (i *Interpreter) VisitWhileStmt(w *parser.While) (interface{}, error) {
	for {
		res, err := i.isTruthy(w.Condition)
		if err != nil {
			return nil, err
		}
		if res.(token.Object).Value_bool == false {
			break
		}
		_, err = i.execute(w.Body)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitVariableExpr(v *parser.Variable) (interface{}, error) {
	return i.Environment.Get(v.Name)
}

func (i *Interpreter) VisitBinaryExpr(b *parser.Binary) (interface{}, error) {
	left, err := i.evaluate(b.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(b.Right)
	if err != nil {
		return nil, err
	}

	switch b.Operator.TokenType {
	case token.GREATER:
		err := i.checkNumberOperands(b.Operator, b.Left, b.Right)
		if err != nil {
			return nil, err
		}
		return token.Object{
			ObjType:    token.BOOL_TYPE,
			Value_bool: left.(token.Object).Value_float > right.(token.Object).Value_float,
		}, nil
	case token.GREATER_EQUAL:
		err := i.checkNumberOperands(b.Operator, b.Left, b.Right)
		if err != nil {
			return nil, err
		}
		return token.Object{
			ObjType:    token.BOOL_TYPE,
			Value_bool: left.(token.Object).Value_float >= right.(token.Object).Value_float,
		}, nil
	case token.LESS:
		err := i.checkNumberOperands(b.Operator, b.Left, b.Right)
		if err != nil {
			return nil, err
		}
		return token.Object{
			ObjType:    token.BOOL_TYPE,
			Value_bool: left.(token.Object).Value_float < right.(token.Object).Value_float,
		}, nil
	case token.LESS_EQUAL:
		err := i.checkNumberOperands(b.Operator, b.Left, b.Right)
		if err != nil {
			return nil, err
		}
		return token.Object{
			ObjType:    token.BOOL_TYPE,
			Value_bool: left.(token.Object).Value_float <= right.(token.Object).Value_float,
		}, nil
	case token.MINUS:
		err := i.checkNumberOperands(b.Operator, b.Left, b.Right)
		if err != nil {
			return nil, err
		}
		return token.Object{
			ObjType:     token.NUMBER_TYPE,
			Value_float: left.(token.Object).Value_float - right.(token.Object).Value_float,
		}, nil
	case token.PLUS:
		leftType := left.(token.Object).ObjType
		rightType := right.(token.Object).ObjType
		if leftType == token.NUMBER_TYPE && rightType == token.NUMBER_TYPE {
			return token.Object{
				ObjType:     token.NUMBER_TYPE,
				Value_float: left.(token.Object).Value_float + right.(token.Object).Value_float,
			}, nil
		} else if leftType == token.STRING_TYPE && rightType == token.STRING_TYPE {
			return token.Object{
				ObjType:   token.STRING_TYPE,
				Value_str: left.(token.Object).Value_str + right.(token.Object).Value_str,
			}, nil
		}

		return nil, &RuntimeError{
			Token: b.Operator,
			Err:   errors.New("Operands must be a number."),
		}
	case token.SLASH:
		err := i.checkNumberOperands(b.Operator, b.Left, b.Right)
		if err != nil {
			return nil, err
		}
		if right.(token.Object).Value_float == 0 {
			return nil, &RuntimeError{
				Token: b.Operator,
				Err:   errors.New("cannot divide by 0"),
			}
		}
		return token.Object{
			ObjType:     token.NUMBER_TYPE,
			Value_float: left.(token.Object).Value_float / right.(token.Object).Value_float,
		}, nil
	case token.STAR:
		err := i.checkNumberOperands(b.Operator, b.Left, b.Right)
		if err != nil {
			return nil, err
		}
		return token.Object{
			ObjType:     token.NUMBER_TYPE,
			Value_float: left.(token.Object).Value_float * right.(token.Object).Value_float,
		}, nil
	case token.BANG_EQUAL:
		val, err := i.isEqual(b.Left, b.Right)
		return token.Object{
			ObjType:    val.(token.Object).ObjType,
			Value_bool: !val.(token.Object).Value_bool,
		}, err
	case token.EQUAL_EQUAL:
		return i.isEqual(b.Left, b.Right)
	}

	return nil, &RuntimeError{
		Token: b.Operator,
		Err:   errors.New("undefined binary token"),
	}
}

func (i *Interpreter) VisitIfStmt(pif *parser.If) (interface{}, error) {
	val, err := i.evaluate(pif.Condition)
	if err != nil {
		return nil, err
	}

	if val.(token.Object).Value_bool {
		_, err = i.execute(pif.ThenBranch)
		if err != nil {
			return nil, err
		}

	} else if pif.ElseBranch != nil {
		_, err = i.execute(pif.ElseBranch)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) isTruthy(e parser.Expr) (interface{}, error) {
	utils.DPrintf("isTruthy -> %+v\n", e)
	v, err := e.Visit(i)
	if err != nil {
		return nil, err
	}
	switch v.(type) {
	case token.Object:
		return i.isObjectTruthy(v.(token.Object))
	}

	return nil, errors.New("undefined type for truthy")
}

func (i *Interpreter) isObjectTruthy(o token.Object) (token.Object, error) {
	switch o.ObjType {
	case token.BOOL_TYPE:
		return o, nil
	case token.STRING_TYPE:
		if len(o.Value_str) > 0 {
			return token.Object{
				ObjType:    token.BOOL_TYPE,
				Value_bool: true,
			}, nil
		} else {
			return token.Object{
				ObjType:    token.BOOL_TYPE,
				Value_bool: false,
			}, nil
		}
	case token.NUMBER_TYPE:
		if o.Value_float != 0 {
			return token.Object{
				ObjType:    token.BOOL_TYPE,
				Value_bool: true,
			}, nil
		} else {
			return token.Object{
				ObjType:    token.BOOL_TYPE,
				Value_bool: false,
			}, nil
		}
	case token.NOT_ASSIGNED_TYPE:
		return token.Object{
			ObjType:    token.BOOL_TYPE,
			Value_bool: false,
		}, nil

	}

	return token.Object{
		ObjType:    token.BOOL_TYPE,
		Value_bool: false,
	}, nil
}

func (i *Interpreter) isEqual(a, b parser.Expr) (interface{}, error) {
	v1, err := a.Visit(i)
	if err != nil {
		return nil, err
	}
	v2, err := b.Visit(i)
	if err != nil {
		return nil, err
	}

	utils.DPrintf("isEqual -> %+v\n%+v\n", v1, v2)
	switch v1.(token.Object).ObjType {
	case token.NUMBER_TYPE:
		if v2.(token.Object).ObjType == token.NUMBER_TYPE {
			return token.Object{
				ObjType:    token.BOOL_TYPE,
				Value_bool: v1.(token.Object).Value_float == v2.(token.Object).Value_float,
			}, nil
		}

	case token.STRING_TYPE:
		if v2.(token.Object).ObjType == token.STRING_TYPE {
			return token.Object{
				ObjType:    token.BOOL_TYPE,
				Value_bool: v1.(token.Object).Value_str == v2.(token.Object).Value_str,
			}, nil
		}
	case token.BOOL_TYPE:
		if v2.(token.Object).ObjType == token.BOOL_TYPE {
			return token.Object{
				ObjType:    token.BOOL_TYPE,
				Value_bool: v1.(token.Object).Value_bool == v2.(token.Object).Value_bool,
			}, nil
		}

	case token.NOT_ASSIGNED_TYPE:
		if v2.(token.Object).ObjType == token.NOT_ASSIGNED_TYPE {
			return token.Object{
				ObjType:    token.BOOL_TYPE,
				Value_bool: true,
			}, nil
		}
	}

	return token.Object{
		ObjType:    token.BOOL_TYPE,
		Value_bool: false,
	}, nil
}

func (i *Interpreter) evaluate(e parser.Expr) (interface{}, error) {
	return e.Visit(i)
}

func (i *Interpreter) checkNumberOperand(operator token.Token, operand parser.Expr) error {
	utils.DPrintf("%+v %+v\n", operator, operand)
	v, err := operand.Visit(i)
	if err != nil {
		return &RuntimeError{
			Token: operator,
			Err:   err,
		}
	}
	if v.(token.Object).ObjType == token.NUMBER_TYPE {
		return nil
	}

	return &RuntimeError{
		Token: operator,
		Err:   errors.New("Operand must be a number."),
	}
}

func (i *Interpreter) checkNumberOperands(operator token.Token, left, right parser.Expr) error {
	utils.DPrintf("%+v %+v %+v\n", operator, left, right)
	type1, err := left.Visit(i)
	if err != nil {
		return &RuntimeError{
			Token: operator,
			Err:   err,
		}
	}

	type2, err := right.Visit(i)
	if err != nil {
		return &RuntimeError{
			Token: operator,
			Err:   err,
		}
	}

	if type1.(token.Object).ObjType == token.NUMBER_TYPE && type2.(token.Object).ObjType == token.NUMBER_TYPE {
		return nil
	}

	return &RuntimeError{
		Token: operator,
		Err:   errors.New("Operands must be a number."),
	}
}

func (i *Interpreter) executeBlock(stmts []parser.Stmt, env *environment.Environment) error {
	previous := i.Environment
	defer func() { i.Environment = previous }()

	i.Environment = env
	for _, stmt := range stmts {
		_, err := i.execute(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}
