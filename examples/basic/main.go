package main

import (
	"fmt"
	"os"

	"github.com/struckchure/gelt"
	"github.com/struckchure/gelt/compiler"
)

func main() {
	componentFile := "examples/basic/counter.gelt"
	code, _ := os.ReadFile(componentFile)

	component, err := compiler.Parse(componentFile, string(code))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(gelt.Prettify(component))

	// res, err := compiler.Walk(*component)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(*res)
}
