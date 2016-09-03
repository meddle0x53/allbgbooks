package models

import (
	"database/sql"
)

type Model interface{}
type Models []Model

type ModelAndFields func() (Model, []interface{})

var ModelAndFiedlsFactories = map[string]ModelAndFields{
	"publishers": func() (Model, []interface{}) {
		var publisher Publisher
		fields := []interface{}{
			&publisher.Id, &publisher.Name, &publisher.Code, &publisher.State,
		}

		return &publisher, fields
	},
	"publisher_aliases": func() (Model, []interface{}) {
		var publisherAlias PublisherAlias
		fields := []interface{}{&publisherAlias.Name}

		return &publisherAlias, fields
	},
	"publisher_contacts": func() (Model, []interface{}) {
		var publisherContact PublisherContact
		fields := []interface{}{&publisherContact.Name}

		return &publisherContact, fields
	},
	"authors": func() (Model, []interface{}) {
		var author Author
		fields := []interface{}{
			&author.Id, &author.Name, &author.Nationality,
		}

		return &author, fields
	},
}

func CreateCollection(rows *sql.Rows, collectionName string) *[]Model {
	defer rows.Close()

	result := []Model{}
	for rows.Next() {
		model, fields := ModelAndFiedlsFactories[collectionName]()

		err := rows.Scan(fields...)
		if err != nil {
			panic(err)
		}

		result = append(result, model)
	}

	return &result
}
