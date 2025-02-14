package types

type Environment struct {
	store map[string]Value
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Value)
	return &Environment{store: s}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(key string) (Value, bool) {
	obj, ok := e.store[key]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(key)
	}
	return obj, ok
}

func (e *Environment) Set(key string, val Value) Value {
	e.store[key] = val
	return val
}
