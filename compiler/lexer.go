package compiler

import "github.com/alecthomas/participle/v2/lexer"

var componentLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "ScriptBlock", Pattern: `(?s)<script\b[^>]*>.*?</script>`},
	{Name: "StyleBlock", Pattern: `(?s)<style\b[^>]*>.*?</style>`},

	{Name: "HTMLComment", Pattern: `(?s)<!--.*?-->`},
	{Name: "JSComment", Pattern: `//.*?$|(?s)/\*.*?\*/`},
	{Name: "CSSComment", Pattern: `(?s)/\*.*?\*/`},

	{Name: "Whitespace", Pattern: `[ \t\r\n]+`},

	{Name: "HTMLBlock", Pattern: `(?s).*`},
})
