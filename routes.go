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
			"Publishers", "GET", "/publishers", actions.CollectionIndexAction,
		),
		routing.NewResourceRoute(
			"Publisher", "GET", "/publishers/{id}", actions.CollectionShowAction,
		),
		routing.NewCollectionRoute(
			"Authors", "GET", "/authors", actions.CollectionIndexAction,
		),
		routing.NewResourceRoute(
			"Author", "GET", "/authors/{id}", actions.CollectionShowAction,
		),
		routing.NewCollectionRoute(
			"Books", "GET", "/books", actions.CollectionIndexAction,
		),
		routing.NewResourceRoute(
			"Book", "GET", "/books/{id}", actions.CollectionShowAction,
		),
	}
}
