package decorators

import (
  "log"
  "net/http"
  "time"
)

func Logger(inner http.Handler, name string) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    now := time.Now()

    inner.ServeHTTP(w, r)

    log.Printf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(now))
  })
}
