package actions

import (
  "net/http"
  "encoding/json"
)

func RenderJson(data interface{}, w http.ResponseWriter, pretty bool) {
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)

  if pretty {
    toRender, err := json.MarshalIndent(data, "", "    ")
    if err != nil { panic(err) }
    w.Write([]byte(toRender))
  } else {
    if err := json.NewEncoder(w).Encode(data); err != nil {
      panic(err)
    }
  }
}
