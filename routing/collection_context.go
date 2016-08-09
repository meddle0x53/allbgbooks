package routing

import "strings"

type CollectionContext struct {
  Context
  PerPage uint64
  Page uint64
  LastPage uint64
  OrderBy string
}

type CollectionAction func(context *CollectionContext)

func NewCollectionRoute(name, method, pattern string, action CollectionAction) Route {
  wrapperAction := func(context Context) {
    collectionContext := context.(*CollectionContext)
    action(collectionContext)
  }

  collectionName := strings.ToLower(name)

  sortAction := Sorting(collectionName, wrapperAction)
  return BasicRoute{
    name, method, pattern, Pagination(collectionName, sortAction),
  }
}

func ToCollectionContext(context Context) *CollectionContext {
  switch context.(type) {
  case *CollectionContext:
    return context.(*CollectionContext)
  default:
    return &CollectionContext{context, 10, 1, 0, "id"}
  }
}
