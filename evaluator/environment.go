package evaluator

import "github.com/siyul-park/miniscript/types"

type Environment struct {
	store map[string]types.Value
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]types.Value)
	return &Environment{store: s}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(key string) (types.Value, bool) {
	obj, ok := e.store[key]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(key)
	}
	return obj, ok
}

func (e *Environment) Set(key string, val types.Value) types.Value {
	e.store[key] = val
	return val
}
