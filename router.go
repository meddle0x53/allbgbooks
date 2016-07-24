package main

import (
  "net/http"
  "allbooks/decorators"
  "github.com/gorilla/mux"
)

type Route struct {
  Name        string
  Method      string
  Pattern     string
  Action      http.HandlerFunc
}

type Routes []Route

func Router(routes Routes) *mux.Router {
  router := mux.NewRouter().StrictSlash(true)

  for _, route := range routes {
    var action http.Handler

    action = route.Action
    action = decorators.Wrap(action, route.Name)

    router.
      Methods(route.Method).
      Path(route.Pattern).
      Name(route.Name).
      Handler(action)
  }

  return router
}
