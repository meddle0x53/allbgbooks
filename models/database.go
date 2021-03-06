package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
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

func Count(context CollectionContext) (uint64, uint64) {
	var result, delta uint64

	filters := context.FilteringValues()
	ignoreCase := context.IgnoreCase()
	name := context.CollectionName()

	query := Filter(sq.Select("count(*)").From(name), filters, ignoreCase)
	query = query.RunWith(GetDB()).PlaceholderFormat(sq.Dollar)

	query.QueryRow().Scan(&result)

	if (result % context.PerPage()) != 0 {
		delta = 1
	}

	return result, (result / context.PerPage()) + delta
}

func GetCollection(context CollectionContext) *sql.Rows {
	collectionName := context.CollectionName()
	selectStatement := strings.Join(CollectionFields[collectionName], ", ")
	offset := (context.Page() - 1) * context.PerPage()

	query := sq.
		Select(selectStatement).
		From(collectionName).
		Limit(context.PerPage()).
		Offset(offset).
		OrderBy(context.OrderBy())

	query = Filter(query, context.FilteringValues(), context.IgnoreCase())
	query = query.RunWith(GetDB()).PlaceholderFormat(sq.Dollar)

	rows, err := query.Query()
	if err != nil {
		panic(err)
	}

	return rows
}

func GetCollectionThroughRelation(field JoinField, id string) *sql.Rows {
	selectStatement := strings.Join(CollectionFields[field.Table], ", ")

	query := sq.
		Select(selectStatement).
		From(field.Table).
		Where(fmt.Sprintf(
			"id IN (SELECT %s FROM %s WHERE %s = ?)", field.TableColumn, field.Type,
			field.Column,
		), id)

	query = query.RunWith(GetDB()).PlaceholderFormat(sq.Dollar)

	rows, err := query.Query()
	if err != nil {
		panic(err)
	}

	return rows
}

func GetResource(collectionName string, context ResourceContext) sq.RowScanner {
	selectFields := CollectionFields[collectionName]
	joinFields := context.JoinFields()
	for _, joinField := range joinFields {
		if joinField.Type == "one" {
			selectFields = append(selectFields, CollectionFields[joinField.Table]...)
		}
	}

	selectStatement := strings.Join(selectFields, ", ")

	query := sq.
		Select(selectStatement).
		From(collectionName)

	id := *context.IdParameter()
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
	"publishers":          []string{"publishers.id", "publishers.name", "code", "state"},
	"publisher_addresses": []string{"town", "main", "phone", "email", "site"},
	"publisher_aliases":   []string{"name"},
	"publisher_contacts":  []string{"name"},
	"languages":           []string{"languages.name"},
	"genres":              []string{"genres.name"},
	"categories":          []string{"categories.name"},
	"authors":             []string{"id", "name", "nationality"},
	"books": []string{
		"books.id", "isbn", "title", "cover", "issue", "description", "copies",
		"price", "publish_date",
	},
}

type IdField struct {
	Pattern string
	Name    string
}

var IdFields = map[string][]IdField{
	"publishers": []IdField{
		IdField{`^\d+$`, "id"}, IdField{`^\d+-\d+-\d+(-\d+)*$`, "code"},
	},
	"authors": []IdField{IdField{`^\d+$`, "id"}},
	"books": []IdField{
		IdField{`^\d+$`, "id"}, IdField{`^\d+-\d+-\d+(-\d+)*$`, "isbn"},
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
	"authors": []FilteringField{
		FilteringField{"name", "LIKE"}, FilteringField{"nationality", "="},
	},
	"books": []FilteringField{
		FilteringField{"isbn", "="}, FilteringField{"title", "LIKE"},
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
		"address":  &JoinField{"id", "publisher_addresses", "publisher_id", "one"},
		"aliases":  &JoinField{"id", "publisher_aliases", "publisher_id", "many"},
		"contacts": &JoinField{"id", "publisher_contacts", "publisher_id", "many"},
	},
	"books": map[string]*JoinField{
		"publisher": &JoinField{"publisher_id", "publishers", "id", "one"},
		"language":  &JoinField{"language_id", "languages", "id", "one"},
		"genre":     &JoinField{"genre_id", "genres", "id", "one"},
		"category":  &JoinField{"category_id", "categories", "id", "one"},
		"authors":   &JoinField{"book_id", "authors", "author_id", "books_authors"},
	},
}
