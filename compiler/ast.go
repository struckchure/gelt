package compiler

type Component struct {
	Script *string `parser:"@ScriptBlock?"`
	Style  *string `parser:"@StyleBlock?"`
	HTML   *string `parser:"@HTMLBlock"`
}
