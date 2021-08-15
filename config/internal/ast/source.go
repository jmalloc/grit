package ast

import "github.com/alecthomas/participle/v2/lexer"

type Source struct {
	Pos        lexer.Position
	Type       string      `parser:"@(\"git\" | \"github\")"`
	NameToken  *string     `parser:"@String?"`
	IsDisabled bool        `parser:"(  @\"disabled\" \";\"?"`
	Parameters []Parameter `parser:" | \"{\" @@* \"}\" )"`
}

func (n Source) Name() string {
	if n.NameToken != nil {
		return unquote(*n.NameToken)
	}

	return n.Type
}

func (n Source) Visit(v Visitor) error {
	return v.VisitSource(n)
}

func (n Source) VisitChildren(v Visitor) error {
	for _, p := range n.Parameters {
		if err := p.Visit(v); err != nil {
			return err
		}
	}

	return nil
}
