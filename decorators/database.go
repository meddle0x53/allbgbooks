package decorators

import(
  "encoding/json"
  "os"
  "fmt"
  "net/http"
  "github.com/gorilla/context"
  "database/sql"
  _ "github.com/lib/pq"
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

func Database(inner http.Handler, name string) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    connection := GetDB()
    context.Set(r, "dbConnection", connection)

    inner.ServeHTTP(w, r)
  })
}
