package routing

import (
	"allbooks/models"
	"strings"
)

type CollectionContext struct {
	Context
	models.CollectionContext
}

type CollectionAction func(context *CollectionContext)

func NewCollectionRoute(name, method, pattern string, action CollectionAction) Route {
	wrapperAction := func(context Context) {
		collectionContext := context.(*CollectionContext)
		action(collectionContext)
	}

	collectionName := strings.ToLower(name)

	sortAction := Sorting(collectionName, wrapperAction)
	paginationAction := Pagination(collectionName, sortAction)
	return BasicRoute{
		name, method, pattern, Filtering(collectionName, paginationAction),
	}
}

func ToCollectionContext(context Context) *CollectionContext {
	switch context.(type) {
	case *CollectionContext:
		return context.(*CollectionContext)
	default:
		return &CollectionContext{context, models.NewCollectionContext()}
	}
}
