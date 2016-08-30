package routing

import (
  "allbooks/models"
)

func Filtering(collection string, action Action) Action {
  return func(context Context) {
    if context.Stop() {
      return
    }

    newContext := ToCollectionContext(context)

    fields := models.FilteringFields[collection]

    ignoreCase := context.Request().Form.Get("ignoreCase")
    if ignoreCase == "true" {
      newContext.IgnoreCase = true
    }

    for _, field := range fields {
      value := context.Request().Form.Get(field.Name)
      if value != "" {
        newContext.FilteringValues = append(
          newContext.FilteringValues, models.FilteringValue{field, value},
        )
      }
    }

    action(newContext)
  }
}
