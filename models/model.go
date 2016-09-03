package models

import (
	"database/sql"
	//	sq "github.com/Masterminds/squirrel"
)

type Model interface {
	Fields() []interface{}
}

type ModelFactory func() Model

var ModelFactories = map[string]ModelFactory{
	"publishers": func() Model {
		var publisher Publisher

		return &publisher
	},
	"publisher_aliases": func() Model {
		var publisherAlias PublisherAlias

		return &publisherAlias
	},
	"publisher_contacts": func() Model {
		var publisherContact PublisherContact

		return &publisherContact
	},
	"authors": func() Model {
		var author Author

		return &author
	},
}

func CreateCollection(rows *sql.Rows, collectionName string) *[]Model {
	defer rows.Close()

	result := []Model{}
	for rows.Next() {
		model := ModelFactories[collectionName]()
		fields := model.Fields()

		err := rows.Scan(fields...)
		if err != nil {
			panic(err)
		}

		result = append(result, model)
	}

	return &result
}

//func CreateResource(context ResourceContext, row sq.RowScanner, collectionName string) Model {
//	model, fields := ModelFactories[collectionName]()
//	joinFields := context.JoinFields()
//
//	for _, joinField := range joinFields {
//		if joinField.Type == "one" {
//			builder := JoinBuilders[joinField.Table]
//
//			fields = builder(model, fields)
//		}
//	}
//
//	return model
//}
