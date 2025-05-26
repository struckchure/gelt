package compiler

import "github.com/alecthomas/participle/v2"

func Parse(filename, code string) (*Component, error) {
	parser, err := participle.Build[Component](
		participle.Lexer(componentLexer),
		participle.Elide("Whitespace", "HTMLComment", "JSComment", "CSSComment"),
	)
	if err != nil {
		return nil, err
	}

	return parser.ParseString(filename, code)
}
