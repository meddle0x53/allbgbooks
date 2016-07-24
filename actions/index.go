package actions

import (
  "net/http"
  "github.com/gorilla/context"
  "allbooks/decorators"
)

func IndexAction(w http.ResponseWriter, r *http.Request) {
  index := struct {
    Message string `json:"message"`
    Publishers []string `json:"publishers"`
  }{
    "Welcome to All BG Books, available endpoints are listed here.",
    []string{
      "/publishers",
      "/publishers/id",
      "/publishers/code",
    },
  }

  params := context.Get(r, "params").(*decorators.Params)
  RenderJson(index, w, params.Pretty)
}
