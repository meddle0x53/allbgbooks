package routing

import (
	"allbooks/models"
	"fmt"
	"github.com/acsellers/inflections"
	"github.com/gorilla/mux"
	"strings"
)

type ResourceContext struct {
	Context
	models.ResourceContext
}

type ResourceAction func(context *ResourceContext)

func handleIdParameter(action Action) Action {
	return func(context Context) {
		if context.Stop() {
			return
		}

		vars := mux.Vars(context.Request())
		newContext := ToResourceContext(context)
		idValue := vars["id"]
		newContext.SetIdParameter(&idValue)

		action(newContext)

		if newContext.IsEmpty() {
			collection := newContext.CollectionName()
			newContext.RespondWithError(
				404, "Not Found", fmt.Sprintf(
					"There is no record in %s identified by %s.", collection, vars["id"],
				),
			)
		}
	}
}

func JoinAction(action Action) Action {
	return func(context Context) {
		if context.Stop() {
			return
		}

		newContext := ToResourceContext(context)
		include := context.Request().Form.Get("include")

		if include != "" {
			includeFields := strings.Split(include, ",")

			for _, includeField := range includeFields {
				trimmed := strings.TrimSpace(includeField)
				toJoin := models.JoinFields[newContext.CollectionName()][trimmed]
				if toJoin != nil {
					newContext.SetJoinFields(append(newContext.JoinFields(), *toJoin))
				} else {
					newContext.RespondWithError(
						422, "Unprocessable Entity",
						fmt.Sprintf("Can't include fields for '%s'.", trimmed),
					)
					return
				}
			}
		}

		action(newContext)
	}
}

func NewResourceRoute(name, method, pattern string, action ResourceAction) Route {
	wrapperAction := func(context Context) {
		resourceContext := context.(*ResourceContext)
		action(resourceContext)
	}

	collectionName := inflections.Underscore(inflections.Pluralize(name))
	setCollectionName := func(action Action) Action {
		return func(context Context) {
			resourceContext := ToResourceContext(context)
			resourceContext.SetCollectionName(collectionName)

			action(resourceContext)
		}
	}

	idAction := handleIdParameter(wrapperAction)
	joinAction := JoinAction(idAction)

	return BasicRoute{name, method, pattern, setCollectionName(joinAction)}
}

func ToResourceContext(context Context) *ResourceContext {
	switch context.(type) {
	case *ResourceContext:
		return context.(*ResourceContext)
	default:
		return &ResourceContext{context, models.NewResourceContext()}
	}
}
