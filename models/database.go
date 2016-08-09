package models

import(
  "encoding/json"
  "os"
  "fmt"
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

func Count(name string, perPage uint64) (uint64, uint64) {
  var result, delta uint64

  query := sq.
    Select("count(*)").
    From(name).
    RunWith(GetDB())

  query.QueryRow().Scan(&result)

  if (result % perPage) != 0 { delta = 1 }

  return result, (result / perPage) + delta
}

var CollectionFields = map[string][]string {
  "publishers": []string{"id", "name", "code", "state",},
}
