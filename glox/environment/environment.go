package environment

import (
	"errors"
	"strconv"

	"github.com/madraceee/interpreters/glox/token"
	"github.com/madraceee/interpreters/glox/utils"
)

type Environment struct {
	Enclosing *Environment
	value     map[string]interface{}
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

func NewEnvironment(env *Environment) *Environment {
	return &Environment{
		Enclosing: env,
		value:     make(map[string]interface{}),
	}
}

func (e *Environment) Define(name string, obj interface{}) {
	e.value[name] = obj
}

func (e *Environment) Get(name token.Token) (interface{}, error) {
	val, ok := e.value[name.Lexeme]
	if !ok {
		if e.Enclosing != nil {
			envclosedVal, err := e.Enclosing.Get(name)
			if envclosedVal == nil {
				return nil, errors.New("Varible '" + name.Lexeme + "' not assigned at line " + strconv.Itoa(name.Line) + ".")
			}
			if err != nil {
				return nil, err
			}
			return envclosedVal, nil
		}

		return nil, &RuntimeError{
			Token: name,
			Err:   errors.New("Undefined variable '" + name.Lexeme + "'"),
		}
	}
	if val == nil {
		return nil, errors.New("Varible '" + name.Lexeme + "' not assigned at line " + strconv.Itoa(name.Line) + ".")
	}

	return val, nil
}

func (e *Environment) Assign(name token.Token, value interface{}) error {
	_, ok := e.Contains(name)
	if ok {
		e.Set(name, value)
		return nil
	}

	if e.Enclosing != nil {
		err := e.Enclosing.Assign(name, value)
		return err
	}

	return &RuntimeError{
		Token: name,
		Err:   errors.New("Undefined variable '" + name.Lexeme + "'."),
	}
}

func (e *Environment) Set(name token.Token, obj interface{}) {
	e.value[name.Lexeme] = obj
}

func (e *Environment) Contains(name token.Token) (interface{}, bool) {
	val, ok := e.value[name.Lexeme]
	return val, ok
}
