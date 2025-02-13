package types

type Kind string

type Value interface {
	Kind() Kind
	Interface() any
}

type Null struct {
}

type Int struct {
	Value int
}

type Float struct {
	Value float64
}

var (
	NULL = &Null{}
)

const (
	KindNull  Kind = "null"
	KindInt   Kind = "int"
	KindFloat Kind = "float"
)

func NewInt(value int) *Int {
	return &Int{Value: value}
}

func NewFloat(value float64) *Float {
	return &Float{Value: value}
}

func (n *Null) Kind() Kind {
	return KindNull
}

func (n *Null) Interface() interface{} {
	return nil
}

func (i *Int) Kind() Kind {
	return KindInt
}

func (i *Int) Interface() any {
	return i.Value
}

func (f *Float) Kind() Kind {
	return KindFloat
}

func (f *Float) Interface() any {
	return f.Value
}
