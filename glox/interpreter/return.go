package interpreter

import "fmt"

type ReturnError struct {
	Value interface{}
}

func (r *ReturnError) Error() string {
	return fmt.Sprint("Return error: %s", r.Value)
}
