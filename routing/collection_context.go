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
	setCollectionName := func(action Action) Action {
		return func(context Context) {
			collectionContext := ToCollectionContext(context)
			collectionContext.SetCollectionName(collectionName)

			action(collectionContext)
		}
	}

	sortAction := Sorting(wrapperAction)
	paginationAction := Pagination(sortAction)
	filteringAction := Filtering(paginationAction)
	return BasicRoute{
		name, method, pattern, setCollectionName(filteringAction),
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
