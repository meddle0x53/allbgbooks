package decorators

import(
  "net/http"
  "github.com/gorilla/context"
)

type Params struct {
  Pretty bool
}

func Parameters(inner http.Handler, name string) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
      panic(err)
    }

    params := &Params{ Pretty: (r.Form.Get("pretty") == "true") }
    context.Set(r, "params", params)

    inner.ServeHTTP(w, r)
  })
}
