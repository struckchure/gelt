package main

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/struckchure/gelt"
)

var iniContent string = `
age = 21
name = "Bob Smith"

[address]
city = "Beverly Hills"
postal_code = 90210
colors = [1, "red", ["green", 9.8]]
env = {"secret": 1, "scheme": "https://"}
`

func main() {
	parser, err := participle.Build[INI](
		participle.Lexer(iniLexer),
		participle.Unquote("String"),
		participle.Union[Value](String{}, Number{}, List{}, Dictionary{}),
		participle.Elide("Whitespace", "SemiColon", "Newline"),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	ini, err := parser.ParseString("ini", iniContent)
	if err != nil {
		fmt.Println(err)
		return
	}

	res, _ := gelt.Prettify(ini)
	fmt.Println(res)
}
