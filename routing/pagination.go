package routing

import (
  "strconv"
  "strings"
  "fmt"
  "allbooks/models"
)

const (
  MaxPerPage = 200
)

func makeLink(
  collection string, page uint64, perPage uint64, rel string,
) string {
  strPage := strconv.FormatUint(page, 10)
  strPerPage := strconv.FormatUint(perPage, 10)

  return `<` + Domain + `/` + collection + `?page=` +
    strPage + `&perPage=` + strPerPage + `>; rel="` + rel + `"`
}

func Pagination(collection string, action Action) Action {
  return func(context Context) {
    if context.Stop() {
      return
    }

    page := ParseIntParam(context, "page", 1)
    perPage := ParseIntParam(context, "perPage", 10)
    count, lastPage := models.Count(collection, perPage)

    if page > lastPage {
      message := fmt.Sprintf(
        "The `page` request param must be lower or equal to %d", lastPage)
      details := fmt.Sprintf(
        "When `perPage` is passed as `%d`, the maximum possible page is " +
        "`%d`, because the number of the results in the DB is `%d`",
        perPage, lastPage, count)
      context.RespondWithError(422, message, details)
      return
    }
    if perPage > MaxPerPage || perPage <= 0 {
      message := fmt.Sprintf(
        "The `perPage` request param must be lower or equal to %d and must " +
        "be greater than zero", MaxPerPage)
      context.RespondWithError(422, message, message)
      return
    }

    newContext := ToCollectionContext(context)
    newContext.PerPage = perPage
    newContext.Page = page
    newContext.LastPage = lastPage

    action(newContext)

    links := make([]string, 0, 4)

    if page < newContext.LastPage {
      links = append(links, makeLink(collection, page + 1, perPage, "next"))
      links = append(
        links, makeLink(collection, newContext.LastPage, perPage, "last"))
    }

    if page > 1 {
      links = append(links, makeLink(collection, 1, perPage, "first"))
      links = append(links, makeLink(collection, page - 1, perPage, "prev"))
    }

    linkHeader := strings.Join(links[:], ", ")

    newContext.SetResponseHeader("Link", linkHeader)
  }
}
