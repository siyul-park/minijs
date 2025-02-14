package ast

type Statement struct {
	Node Node
}

func NewStatement(node Node) *Statement {
	return &Statement{Node: node}
}

func (n *Statement) String() string {
	return n.Node.String() + ";"
}
