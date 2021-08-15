package ast

// Visitable is an interface for AST nodes that can be visited.
type Visitable interface {
	Visit(Visitor) error
}

// Visitor is an interface for visiting AST nodes.
type Visitor interface {
	VisitSource(Source) error
	VisitParameter(Parameter) error
}
