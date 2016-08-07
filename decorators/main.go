package decorators

import (
  "net/http"
)


func Decorators() [](func(http.Handler, string) http.Handler) {
  return [](func(http.Handler, string) http.Handler) {
    Logger,
  }
}

func Wrap(handler http.Handler, name string) http.Handler {
  result := handler

  for _, decorator := range Decorators() {
    result = decorator(result, name)
  }

  return result
}
