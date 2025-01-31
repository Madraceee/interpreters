package interpreter

import (
	"time"

	"github.com/madraceee/interpreters/glox/environment"
	"github.com/madraceee/interpreters/glox/parser"
)

type LoxCallable interface {
	Arity() int
	Call(i *Interpreter, arguments []interface{}) (interface{}, error)
}

// In-built clock function
type clock struct{}

func (c *clock) Arity() int {
	return 0
}

func (c *clock) Call(i *Interpreter, _obj []interface{}) (interface{}, error) {
	return time.Now(), nil
}

// User-defined Function
type LoxFunction struct {
	Declaration *parser.Function
}

func NewLoxFunction(declaration *parser.Function) *LoxFunction {
	return &LoxFunction{
		Declaration: declaration,
	}
}

func (lf *LoxFunction) Call(i *Interpreter, arguments []interface{}) (interface{}, error) {
	environment := environment.NewEnvironment(i.globals)
	for i, dec := range lf.Declaration.Params {
		environment.Define(dec.Lexeme, arguments[i])
	}

	err := i.executeBlock(lf.Declaration.Body, environment)
	return nil, err
}

func (lf *LoxFunction) Arity() int {
	return len(lf.Declaration.Params)
}

func (lf *LoxFunction) String() string {
	return "<fn " + lf.Declaration.Name.Lexeme + " >"
}
