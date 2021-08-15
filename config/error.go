package config

import (
	"errors"
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/jmalloc/grit/config/internal/ast"
)

// newError returns an error that indicates a problem with a configuration file.
func newError(pos lexer.Position, format string, args ...interface{}) error {
	return errors.New(
		pos.String() + " " + fmt.Sprintf(format, args...),
	)
}

// missingParameter returns an error that indicates a required parameter was
// absent.
func missingParameterInSource(pos lexer.Position, source, expectedKey string) error {
	return newError(
		pos,
		`missing required "%s" parameter in "%s" source`,
		expectedKey,
		source,
	)
}

// unrecognizedParameter returns an error that indicates an unrecognized
// parameter has been encountered.
func unrecognizedParameter(p ast.Parameter) error {
	return newError(
		p.Pos,
		`unrecognized "%s" parameter`,
		p.Key,
	)
}

// unrecognizedParameterInSource returns an error that indicates an unrecognized
// parameter has been encountered.
func unrecognizedParameterInSource(p ast.Parameter, source string) error {
	return newError(
		p.Pos,
		`unrecognized "%s" parameter in "%s" source`,
		p.Key,
		source,
	)
}

// invalidParameterType returns an error that indicates a parameter has a value
// of an unexpected type.
func invalidParameterType(p ast.Parameter, expectedType string) error {
	return newError(
		p.Value.Pos,
		`invalid type for "%s" parameter, expected %s`,
		p.Key,
		expectedType,
	)
}

// invalidParameterValue returns an error that indicates a parameter has an
// unexpected value despite being the correct type.
func invalidParameterValue(p ast.Parameter, format string, args ...interface{}) error {
	return newError(
		p.Value.Pos,
		`invalid value for "%s" parameter, %s`,
		p.Key,
		fmt.Sprintf(format, args...),
	)

}
