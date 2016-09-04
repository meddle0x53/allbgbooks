package routing

import (
	"strings"
)

func getSortQuery(sortParam string) string {
	sortData := strings.Split(sortParam, ":")

	var result = ""

	if len(sortData) > 0 && sortData[0] != "" {
		result = sortData[0]
	}

	if len(sortData) > 1 && sortData[1] != "" {
		if sortData[1] == "d" {
			result = result + " DESC"
		} else if sortData[1] == "a" {
			result = result + " ASC"
		} else {
			// Error
		}
	}

	return result
}

func setOrderBy(context *CollectionContext, sortParam string) {
	var orderBy = make([]string, 0, 0)

	for _, sortPart := range strings.Split(sortParam, ",") {
		trimmedSortPart := strings.TrimSpace(sortPart)
		query := getSortQuery(trimmedSortPart)
		if query != "" {
			orderBy = append(orderBy, query)
		}
	}

	if len(orderBy) > 0 {
		context.SetOrderBy(strings.Join(orderBy[:], ","))
	}
}

func Sorting(action Action) Action {
	return func(context Context) {
		if context.Stop() {
			return
		}

		newContext := ToCollectionContext(context)

		sortParam := context.Request().Form.Get("sort")
		if sortParam != "" {
			setOrderBy(newContext, sortParam)
		}

		action(newContext)
	}
}
