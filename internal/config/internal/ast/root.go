package ast

type Root struct {
	Statements []Statement `parser:"@@*"`
}

func (n Root) Visit(v Visitor) error {
	for _, stmt := range n.Statements {
		if err := stmt.Visit(v); err != nil {
			return err
		}
	}

	return nil
}
