package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/stdlib"
	"os"
	"regexp"
	"strings"
)

type DBCredentials struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	DBName   string `json:"dbname"`
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

	db, err := sql.Open("pgx", connString)

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

	if (result % perPage) != 0 {
		delta = 1
	}

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
	if err != nil {
		panic(err)
	}

	return rows
}

func GetResource(collectionName, id string, joinFields []JoinField) sq.RowScanner {
	selectFields := CollectionFields[collectionName]
	for _, joinField := range joinFields {
		if joinField.Type == "one" {
			selectFields = append(selectFields, CollectionFields[joinField.Table]...)
		}
	}

	selectStatement := strings.Join(selectFields, ", ")

	query := sq.
		Select(selectStatement).
		From(collectionName)

	idFields := make([]string, 0, len(IdFields[collectionName]))
	for _, idField := range IdFields[collectionName] {
		match, _ := regexp.MatchString(idField.Pattern, id)
		if match {
			idFields = append(idFields, idField.Name)
		}
	}

	orExpression := make(sq.Or, 0, len(idFields))
	for _, idField := range idFields {
		idColumn := fmt.Sprintf("%s.%s", collectionName, idField)
		orExpression = append(
			orExpression, sq.Expr(fmt.Sprintf("%s = ?", idColumn), id),
		)
	}
	query = query.Where(orExpression)

	for _, joinField := range joinFields {
		if joinField.Type == "one" {
			joinStatement := fmt.Sprintf(
				"%s ON %s.%s = %s.%s", joinField.Table, joinField.Table,
				joinField.TableColumn, collectionName, joinField.Column,
			)
			query = query.Join(joinStatement)
		}
	}

	query = query.RunWith(GetDB()).PlaceholderFormat(sq.Dollar)

	return query.QueryRow()
}

var CollectionFields = map[string][]string{
	"publishers":          []string{"publishers.id", "name", "code", "state"},
	"publisher_addresses": []string{"town", "main", "phone", "email", "site"},
	"publisher_aliases":   []string{"name"},
}

type IdField struct {
	Pattern string
	Name    string
}

var IdFields = map[string][]IdField{
	"publishers": []IdField{
		IdField{`^\d+$`, "id"}, IdField{`^\d+-\d+-\d+(-\d+)*$`, "code"},
	},
}

type FilteringField struct {
	Name      string
	QueryType string
}

type FilteringValue struct {
	FilteringField
	Value string
}

var FilteringFields = map[string][]FilteringField{
	"publishers": []FilteringField{
		FilteringField{"name", "LIKE"}, FilteringField{"code", "="},
		FilteringField{"state", "="}, FilteringField{"id", "="},
	},
}

type JoinField struct {
	Column      string
	Table       string
	TableColumn string
	Type        string
}

var JoinFields = map[string]map[string]*JoinField{
	"publishers": map[string]*JoinField{
		"address": &JoinField{"id", "publisher_addresses", "publisher_id", "one"},
		"aliases": &JoinField{"id", "publisher_aliases", "publisher_id", "many"},
	},
}
