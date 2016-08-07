package routing

import (
  "strconv"
  "strings"
  "allbooks/models"
)

func ParseIntParam(ctx Context, name string, defaultVal uint64) uint64 {
  val, err := strconv.ParseUint(ctx.Request().Form.Get(name), 10, 64)
  if err != nil { val = defaultVal }

  return val
}

func makeLink(collection string, page uint64, perPage uint64, rel string) string {
  strPage := strconv.FormatUint(page, 10)
  strPerPage := strconv.FormatUint(perPage, 10)

  return `<` + Domain + `/` + collection + `?page=` +
    strPage + `&perPage=` + strPerPage + `>; rel="` + rel + `"`
}

func Pagination(collection string, action Action) Action {
  return func(context Context) {
    page := ParseIntParam(context, "page", 1)
    perPage := ParseIntParam(context, "perPage", 10)

    newContext := ToCollectionContext(context)
    newContext.PerPage = perPage
    newContext.Page = page
    newContext.LastPage = models.Count(collection, perPage)

    action(newContext)

    links := make([]string, 0, 4)

    if page < newContext.LastPage {
      links = append(links, makeLink(collection, page + 1, perPage, "next"))
      links = append(links, makeLink(collection, newContext.LastPage, perPage, "last"))
    }

    if page > 1 {
      links = append(links, makeLink(collection, 1, perPage, "first"))
      links = append(links, makeLink(collection, page - 1, perPage, "prev"))
    }

    linkHeader := strings.Join(links[:], ", ")

    newContext.SetResponseHeader("Link", linkHeader)
  }
}
