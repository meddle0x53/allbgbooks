package routing

import (
  "fmt"
  "strings"
)

func Sorting(collection string, action Action) Action {
  return func(context Context) {
    if context.Stop() {
      return
    }

    sortParam := context.Request().Form.Get("sort")
    sortData := strings.Split(sortParam, ":")

    newContext := ToCollectionContext(context)
    if len(sortData) > 0 && sortData[0] != "" {
      newContext.OrderBy = sortData[0]
    }

    if len(sortData) > 1 && sortData[1] != "" {
      if sortData[1] == "d" {
        newContext.OrderBy = newContext.OrderBy + " DESC"
      } else if sortData[1] == "a" {
        newContext.OrderBy = newContext.OrderBy + " ASC"
      } else {
        // Error
      }
    }

    fmt.Println(newContext)

    action(newContext)
  }
}
