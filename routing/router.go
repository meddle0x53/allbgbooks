package routing

import (
  "net/http"
  "encoding/json"
  "allbooks/decorators"
  "github.com/gorilla/mux"
)

const Domain = "http://0.0.0.0:8081"

type Action func(Context)

type Route interface {
  Name() string
  Method() string
  Pattern() string
  Action() Action
  AroundActions() []Action
}

type BasicRoute struct {
  RouteName        string
  RouteMethod      string
  RoutePattern     string
  RouteAction      Action
}

func (route BasicRoute) Name() string {
  return route.RouteName
}

func (route BasicRoute) Method() string {
  return route.RouteMethod
}

func (route BasicRoute) Pattern() string {
  return route.RoutePattern
}

func (route BasicRoute) Action() Action {
  return route.RouteAction
}

func (route BasicRoute) AroundActions() []Action {
  return make([]Action, 0, 0)
}

type CompositeRoute struct {
  Route
  RouteAroundActions []Action
}

func (route CompositeRoute) AroundActions() []Action {
  return route.RouteAroundActions
}

type Routes []Route

func wrapAction(action Action) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
      panic(err)
    }

    context := &BasicContext{
      request: r, response: w, responseHeaders: make(map[string]string),
    }

    action(context)
    data := context.ResponseData()

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    for name, val := range context.ResponseHeaders() {
      w.Header().Set(name, val)
    }

    w.WriteHeader(context.Status())

    pretty := r.Form.Get("pretty") == "true"
    if pretty {
      toRender, err := json.MarshalIndent(data, "", "  ")
      if err != nil { panic(err) }
      w.Write([]byte(toRender))
    } else {
      if err := json.NewEncoder(w).Encode(data); err != nil {
        panic(err)
      }
    }
  })
}

func Router(routes Routes) *mux.Router {
  router := mux.NewRouter().StrictSlash(true)

  for _, route := range routes {
    var action http.Handler

    action = wrapAction(route.Action())
    action = decorators.Wrap(action, route.Name())

    router.
      Methods(route.Method()).
      Path(route.Pattern()).
      Name(route.Name()).
      Handler(action)
  }

  return router
}
