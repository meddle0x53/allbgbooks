package routing

import (
	"allbooks/models"
)

func Filtering(action Action) Action {
	return func(context Context) {
		if context.Stop() {
			return
		}

		newContext := ToCollectionContext(context)

		fields := models.FilteringFields[newContext.CollectionName()]

		ignoreCase := context.Request().Form.Get("ignoreCase")
		newContext.SetIgnoreCase(ignoreCase == "true")

		for _, field := range fields {
			value := context.Request().Form.Get(field.Name)
			if value != "" {
				newContext.SetFilteringValues(append(
					newContext.FilteringValues(), models.FilteringValue{field, value},
				))
			}
		}

		action(newContext)
	}
}
