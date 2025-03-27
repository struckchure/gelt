package main

import "github.com/struckchure/gelt"

func main() {
	app := gelt.NewServer(PageRegistry)
	app.Start(8080)
}
