package models

import(
  "encoding/json"
  "os"
  "fmt"
  "strings"
  "database/sql"
  _ "github.com/lib/pq"
  sq "github.com/Masterminds/squirrel"
)

type DBCredentials struct {
  User string `json:"user"`
  Password string `json:"password"`
  Host string `json:"host"`
  DBName string `json:"dbname"`
}

type DBConfiguration struct {
  Production DBCredentials `json:"production"`
}

var dbConfiguration *DBConfiguration
func GetDBConfiguration() *DBConfiguration {
  if dbConfiguration != nil {
    return dbConfiguration
  }

  file, _ := os.Open("config/database.json")
  decoder := json.NewDecoder(file)
  dbConfiguration = &DBConfiguration{}
  err := decoder.Decode(dbConfiguration)

  if err != nil {
    panic(err)
  }

  return dbConfiguration
}

var db *sql.DB
func GetDB() *sql.DB {
  if db != nil {
    return db
  }

  config := GetDBConfiguration().Production
  connString := fmt.Sprintf(
    "postgres://%s:%s@%s/%s",
    config.User, config.Password, config.Host, config.DBName,
  )

  db, err := sql.Open("postgres", connString)

  if err != nil {
    panic(err)
  }

  return db
}


func Filter(
  query sq.SelectBuilder, filters []FilteringValue, ignoreCase bool,
) sq.SelectBuilder {
  for _, filter := range filters {
    if filter.QueryType == "LIKE" {
      filter.Value = "%" + filter.Value + "%"
    }

    placeholder := "%s %s ?"
    if ignoreCase {
      placeholder = "LOWER(%s) %s LOWER(?)"
    }

    query = query.Where(
      fmt.Sprintf(placeholder, filter.Name, filter.QueryType), filter.Value,
    )
  }

  return query
}

func Count(
  name string, perPage uint64, filters []FilteringValue, ignoreCase bool,
) (uint64, uint64) {
  var result, delta uint64

  query := Filter(sq.Select("count(*)").From(name), filters, ignoreCase)
  query = query.RunWith(GetDB()).PlaceholderFormat(sq.Dollar)

  query.QueryRow().Scan(&result)

  if (result % perPage) != 0 { delta = 1 }

  return result, (result / perPage) + delta
}

func GetCollection(
  collectionName string, page uint64, perPage uint64,
  orderBy string, filters []FilteringValue, ignoreCase bool,
) *sql.Rows {
  selectStatement := strings.Join(CollectionFields[collectionName], ", ")
  offset := (page - 1) * perPage

  query := sq.
    Select(selectStatement).
    From(collectionName).
    Limit(perPage).
    Offset(offset).
    OrderBy(orderBy)

  query = Filter(query, filters, ignoreCase)
  query = query.RunWith(GetDB()).PlaceholderFormat(sq.Dollar)

  rows, err := query.Query()
  if err != nil { panic(err) }

  return rows
}

var CollectionFields = map[string][]string {
  "publishers": []string{"id", "name", "code", "state",},
}

type FilteringField struct {
  Name string
  QueryType string
}

type FilteringValue struct {
  FilteringField
  Value string
}

var FilteringFields = map[string][]FilteringField {
  "publishers": []FilteringField{
    FilteringField{"name", "LIKE"}, FilteringField{"code", "="},
    FilteringField{"state", "="},
  },
}
