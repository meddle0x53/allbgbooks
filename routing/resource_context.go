package routing

import (
	"fmt"
	"github.com/acsellers/inflections"
	"github.com/gorilla/mux"
)

type ResourceContext struct {
	Context
	IncludeRelations bool
	IdParameter      string
	IsEmpty          bool
}

type ResourceAction func(context *ResourceContext)

func handleIdParameter(collection string, action Action) Action {
	return func(context Context) {
		if context.Stop() {
			return
		}

		vars := mux.Vars(context.Request())
		newContext := ToResourceContext(context)
		newContext.IdParameter = vars["id"]

		action(newContext)

		if newContext.IsEmpty {
			newContext.RespondWithError(
				404, "Not Found", fmt.Sprintf(
					"There is no record in %s identified by %s.", collection, vars["id"],
				),
			)
		}
	}
}

func NewResourceRoute(name, method, pattern string, action ResourceAction) Route {
	wrapperAction := func(context Context) {
		resourceContext := context.(*ResourceContext)
		action(resourceContext)
	}

	collectionName := inflections.Underscore(inflections.Pluralize(name))
	idAction := handleIdParameter(collectionName, wrapperAction)

	return BasicRoute{name, method, pattern, idAction}
}

func ToResourceContext(context Context) *ResourceContext {
	switch context.(type) {
	case *ResourceContext:
		return context.(*ResourceContext)
	default:
		return &ResourceContext{context, false, "", false}
	}
}
