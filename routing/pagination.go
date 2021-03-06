package routing

import (
	"allbooks/models"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const (
	MaxPerPage = 200
)

func makeLink(
	context *CollectionContext, page uint64, perPage uint64, rel string,
) string {
	strPage := strconv.FormatUint(page, 10)
	strPerPage := strconv.FormatUint(perPage, 10)

	Url, err := url.Parse(Domain)
	if err != nil {
		panic(err)
	}
	Url.Path += context.CollectionName()
	parameters := url.Values{}

	for key, values := range context.Request().Form {
		if key != "page" && key != "perPage" {
			for _, value := range values {
				parameters.Add(key, value)
			}
		}
	}
	parameters.Add("page", strPage)
	parameters.Add("perPage", strPerPage)
	Url.RawQuery = parameters.Encode()

	return `<` + Url.String() + `>; rel="` + rel + `"`
}

func Pagination(action Action) Action {
	return func(context Context) {
		if context.Stop() {
			return
		}

		newContext := ToCollectionContext(context)

		page := ParseIntParam(context, "page", 1)
		perPage := ParseIntParam(context, "perPage", 10)

		newContext.SetPerPage(perPage)
		newContext.SetPage(page)

		count, lastPage := models.Count(newContext)

		if count == 0 {
			context.RespondWithError(
				400, "Not Found", "Nothing was found for this request.",
			)
			return
		}

		if page > lastPage {
			message := fmt.Sprintf(
				"The `page` request param must be lower or equal to %d", lastPage)
			details := fmt.Sprintf(
				"When `perPage` is passed as `%d`, the maximum possible page is "+
					"`%d`, because the number of the results in the DB is `%d`",
				perPage, lastPage, count)
			context.RespondWithError(422, message, details)
			return
		}
		if perPage > MaxPerPage || perPage <= 0 {
			message := fmt.Sprintf(
				"The `perPage` request param must be lower or equal to %d and must "+
					"be greater than zero", MaxPerPage)
			context.RespondWithError(422, message, message)
			return
		}

		action(newContext)

		links := make([]string, 0, 4)

		if page < lastPage {
			links = append(links, makeLink(newContext, page+1, perPage, "next"))
			links = append(
				links, makeLink(newContext, lastPage, perPage, "last"))
		}

		if page > 1 {
			links = append(links, makeLink(newContext, 1, perPage, "first"))
			links = append(links, makeLink(newContext, page-1, perPage, "prev"))
		}

		linkHeader := strings.Join(links[:], ", ")

		newContext.SetResponseHeader("Link", linkHeader)
	}
}
