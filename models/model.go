package models

import (
	"database/sql"
	"strconv"
)

type Model interface {
	Fields() []interface{}
	AppendRelationFields(JoinField, []interface{}) (Model, []interface{})
	Identifier() string
	AppendCollection(string, []FilteringValue) *[]Model
	IsEmpty() bool
	SetRelation(string, *[]Model)
}

type BaseModel struct{}

func (BaseModel) Fields() []interface{} {
	return make([]interface{}, 0)
}

func (m *BaseModel) AppendRelationFields(f JoinField, fields []interface{}) (Model, []interface{}) {
	if f.Type == "one" {
		model := ModelFactories[f.Table]()

		return model, append(fields, model.Fields()...)
	}

	return m, fields
}

func (BaseModel) SetRelation(name string, models *[]Model) {}

func (BaseModel) Identifier() string {
	return ""
}

func (m *BaseModel) IsEmpty() bool {
	return true
}

func (BaseModel) AppendCollection(name string, filters []FilteringValue) *[]Model {
	context := &BaseCollectionContext{100, 1, "id", filters, false}
	rows := GetCollection(name, context)
	return CreateCollection(rows, name)
}

type ModelWithId struct {
	BaseModel
	Id int `json:"id"`
}

func (model *ModelWithId) Identifier() string {
	return strconv.Itoa(model.Id)
}

func (model *ModelWithId) IsEmpty() bool {
	return model.Id <= 0
}

type ModelFactory func() Model

var ModelFactories = map[string]ModelFactory{
	"publishers": func() Model {
		var publisher Publisher

		return &publisher
	},
	"publisher_addresses": func() Model {
		var address PublisherAddress

		return &address
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

func CreateResource(context ResourceContext, collectionName string) Model {
	model := ModelFactories[collectionName]()
	fields := model.Fields()
	joinFields := context.JoinFields()

	for _, joinField := range joinFields {
		if joinField.Type == "one" {
			var relation Model
			relation, fields = model.AppendRelationFields(joinField, fields)
			model.SetRelation(joinField.Table, &[]Model{relation})
		}
	}

	GetResource(collectionName, context).Scan(fields...)
	context.SetIsEmpty(model.IsEmpty())

	for _, joinField := range joinFields {
		if joinField.Type == "many" {
			filters := []FilteringValue{
				FilteringValue{
					FilteringField{joinField.TableColumn, "="}, model.Identifier(),
				},
			}

			collection := model.AppendCollection(joinField.Table, filters)
			model.SetRelation(joinField.Table, collection)
		}
	}

	return model
}
