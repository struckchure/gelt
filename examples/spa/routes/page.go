package index

type Page struct{}

func say_hi() {
	// dom.Alert("Hi")
}

func (Page) Methods() map[string]any {
	return map[string]any{
		"say_hi": say_hi,
	}
}
