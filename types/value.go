package types

type Kind string

type Value interface {
	Kind()
	Interface() any
}
