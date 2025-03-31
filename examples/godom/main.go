package main

//go:generate gopherjs build main.go -o main.js -m

import (
	"fmt"

	dom "github.com/siongui/godom"
)

func main() {
	button := dom.Document.CreateElement("button")
	button.SetTextContent("Click Me")
	button.AddEventListener("click", func(e dom.Event) {
		fmt.Println("click")
		dom.Alert("Hello, World!")
	})

	root := dom.Document.QuerySelector("#root")
	root.AppendChild(button)
}
