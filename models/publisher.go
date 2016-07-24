package models

import (
  "database/sql"
)

type Publisher struct {
  Id int `json:"id"`
  Name string `json:"name"`
  Code string `json:"code"`
  State string `json:"state"`
}

type Publishers []Publisher

func GetPublishers(db *sql.DB) *Publishers {
  rows, err := db.Query(`SELECT id, name, code, state FROM publishers LIMIT 10`)

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
