package main

type INI struct {
	Properties []*Property `parser:"@@*"`
	Sections   []*Section  `parser:"@@*"`
}

type Section struct {
	Identifier string      `parser:"OpenSB @Ident CloseSB"`
	Properties []*Property `parser:"@@*"`
}

type Property struct {
	Key   string `parser:"@Ident Equals"`
	Value Value  `parser:"@@"`
}

type Value interface{ value() }

type String struct {
	String string `parser:"@String"`
}

func (String) value() {}

type Number struct {
	Number float64 `parser:"@(Float | Int)"`
}

func (Number) value() {}

type List struct {
	List []Value `parser:"OpenSB (@@ (Comma @@)*)? CloseSB"`
}

func (List) value() {}

type Dictionary struct {
	Pairs []KV `parser:"OpenCB (@@ (Comma @@)*)? CloseCB"`
}

func (Dictionary) value() {}

type KV struct {
	Key   string `parser:"@String Colon"`
	Value Value  `parser:"@@"`
}

func (KV) value() {}
