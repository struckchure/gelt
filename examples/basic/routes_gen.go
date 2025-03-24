package main

import (
	index "github.com/struckchure/gelt/examples/basic/routes"
	post_list "github.com/struckchure/gelt/examples/basic/routes/posts"
	post_details "github.com/struckchure/gelt/examples/basic/routes/posts/_id"
	post_analytics "github.com/struckchure/gelt/examples/basic/routes/posts/_id/analytics"
	post_create "github.com/struckchure/gelt/examples/basic/routes/posts/create"
	profile "github.com/struckchure/gelt/examples/basic/routes/profile"
)

var PageRegistry = map[string]any{
	"index":          index.Page{},
	"profile":        profile.Page{},
	"post_list":      post_list.Page{},
	"post_details":   post_details.Page{},
	"post_analytics": post_analytics.Page{},
	"post_create":    post_create.Page{},
}
