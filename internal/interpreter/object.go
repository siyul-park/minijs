package interpreter

type Object struct {
	prototype  *Object
	properties map[Value]Value
}

func NewObject(prototype *Object) *Object {
	return &Object{
		prototype:  prototype,
		properties: make(map[Value]Value),
	}
}

func (o *Object) Prototype() *Object {
	return o.prototype
}

func (o *Object) Get(key Value) (Value, bool) {
	value, ok := o.properties[key]
	if !ok && o.prototype != nil {
		return o.prototype.Get(key)
	}
	return value, ok
}

func (o *Object) Set(key, val Value) {
	o.properties[key] = val
}

func (o *Object) Type() Type {
	return OBJECT
}

func (o *Object) Interface() any {
	return nil
}
