package ast

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Parameter struct {
	Pos   lexer.Position
	Key   string `parser:"@Ident \"=\""`
	Value Value  `parser:"@@ \";\""`
}

func (n Parameter) Visit(v Visitor) error {
	return v.VisitParameter(n)
}

type Value struct {
	Pos lexer.Position
	// BoolToken   *bool   `parser:"  @Bool"`
	StringToken *string `parser:"@String"`
}

// func (v *Value) AsBool() (bool, bool) {
// 	if v.BoolToken != nil {
// 		return *v.BoolToken, true
// 	}

// 	return false, false
// }

func (v *Value) AsString() (string, bool) {
	if v.StringToken != nil {
		return unquote(*v.StringToken), true
	}

	return "", false
}
