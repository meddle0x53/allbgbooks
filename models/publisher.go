package models

import sq "github.com/Masterminds/squirrel"

type Publisher struct {
  Id int `json:"id"`
  Name string `json:"name"`
  Code string `json:"code"`
  State string `json:"state"`
}

type Publishers []Publisher

func GetPublishers(page uint64, perPage uint64) *Publishers {
  offset := (page - 1) * perPage

  query := sq.
    Select("id, name, code, state").
    From("publishers").
    Limit(perPage).
    Offset(offset).
    RunWith(GetDB()).
    PlaceholderFormat(sq.Dollar)


  rows, err := query.Query()

  if err != nil { panic(err) }
  defer rows.Close()

  result := Publishers{}
  for rows.Next() {
    var publisher Publisher

    err := rows.Scan(
      &publisher.Id, &publisher.Name, &publisher.Code, &publisher.State)
    if err != nil { panic(err) }

    result = append(result, publisher)
  }

  return &result
}
