package main

import (
	"allbooks/actions"
	"allbooks/routing"
)

func AppRoutes() routing.Routes {
	return routing.Routes{
		routing.BasicRoute{
			"Index", "GET", "/", actions.IndexAction,
		},
		routing.NewCollectionRoute(
			"Publishers", "GET", "/publishers", actions.PublishersIndexAction,
		),
		routing.NewResourceRoute(
			"Publisher", "GET", "/publishers/{id}", actions.PublishersShowAction,
		),
		routing.NewCollectionRoute(
			"Authors", "GET", "/authors", actions.AuthorsIndexAction,
		),
		routing.NewResourceRoute(
			"Author", "GET", "/authors/{id}", actions.AuthorsShowAction,
		),
	}
}
