package ast

type Statement struct {
	Parameter *Parameter `parser:"  @@"`
	Source    *Source    `parser:"| @@"`
}

func (n Statement) Visit(v Visitor) error {
	if n.Parameter != nil {
		return v.VisitParameter(*n.Parameter)
	}

	if n.Source != nil {
		return v.VisitSource(*n.Source)
	}

	panic("unrecognised statement")
}
