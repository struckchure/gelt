package main

import "github.com/alecthomas/participle/v2/lexer"

var iniLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Ident", Pattern: `[a-zA-Z_]\w*`},
	{Name: "String", Pattern: `"(\\"|[^"])*"`},

	{Name: "Float", Pattern: `[-+]?(?:\d+\.\d*|\.\d+)(?:[eE][-+]?\d+)?`},
	{Name: "Int", Pattern: `[-+]?\d+`},

	{Name: "Equals", Pattern: `=`},
	{Name: "Comma", Pattern: `,`},
	{Name: "OpenSB", Pattern: `\[`},
	{Name: "CloseSB", Pattern: `\]`},
	{Name: "OpenCB", Pattern: `\{`},
	{Name: "CloseCB", Pattern: `\}`},

	{Name: "SemiColon", Pattern: `;`},
	{Name: "Colon", Pattern: `:`},

	{Name: "Whitespace", Pattern: `[ \t]+`},
	{Name: "Newline", Pattern: `\r?\n`},
})
