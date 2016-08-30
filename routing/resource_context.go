package routing

type ResourceContext struct {
  Context
  IncludeRelations bool
}

type ResourceAction func(context *ResourceContext)

func NewResourceRoute(name, method, pattern string, action ResourceAction) Route {
  wrapperAction := func(context Context) {
    resourceContext := context.(*ResourceContext)
    action(resourceContext)
  }

  tmp := func(collection string, action Action) Action {
    return func(context Context) {
      newContext := ToResourceContext(context)

      action(newContext)
    }
  }

  return BasicRoute{name, method, pattern, tmp("", wrapperAction)}
}

func ToResourceContext(context Context) *ResourceContext {
  switch context.(type) {
  case *ResourceContext:
    return context.(*ResourceContext)
  default:
    return &ResourceContext{context, false}
  }
}

