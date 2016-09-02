package models

import (
	"strconv"
)

type Publisher struct {
	Id      int               `json:"id"`
	Name    string            `json:"name"`
	Code    string            `json:"code"`
	State   string            `json:"state"`
	Address *PublisherAddress `json:"address,omitempty"`
	Aliases *[]PublisherAlias `json:"aliases,omitempty"`
}

type Publishers []Publisher

func GetPublishers(
	page uint64, perPage uint64, orderBy string,
	filters []FilteringValue, ignoreCase bool,
) *Publishers {
	rows := GetCollection(
		"publishers", page, perPage, orderBy, filters, ignoreCase,
	)
	defer rows.Close()

	result := Publishers{}
	for rows.Next() {
		var publisher Publisher

		err := rows.Scan(
			&publisher.Id, &publisher.Name, &publisher.Code, &publisher.State,
		)
		if err != nil {
			panic(err)
		}

		result = append(result, publisher)
	}

	return &result
}

func GetPublisherById(id string, joinFields []JoinField) *Publisher {
	var publisher Publisher

	fields := []interface{}{
		&publisher.Id, &publisher.Name, &publisher.Code, &publisher.State,
	}

	for _, joinField := range joinFields {
		if joinField.Type == "one" {
			builder := JoinBuilders[joinField.Table]

			fields = builder(&publisher, fields)
		}
	}

	GetResource("publishers", id, joinFields).Scan(fields...)

	for _, joinField := range joinFields {
		if joinField.Type == "many" {
			filters := []FilteringValue{
				FilteringValue{
					FilteringField{joinField.TableColumn, "="},
					strconv.Itoa(publisher.Id),
				},
			}

			builder := CollectionBuilders[joinField.Table]
			builder(&publisher, filters)
		}
	}

	return &publisher
}

func buildPublisherAliases(p *Publisher, filters []FilteringValue) {
	rows := GetCollection("publisher_aliases", 1, 100, "id", filters, false)
	defer rows.Close()

	result := make([]PublisherAlias, 0)
	for rows.Next() {
		var publisher_alias PublisherAlias

		err := rows.Scan(&publisher_alias.Name)
		if err != nil {
			panic(err)
		}

		result = append(result, publisher_alias)
	}

	p.Aliases = &result
}

type JoinBuilder func(*Publisher, []interface{}) []interface{}

func buildPublisherAddress(p *Publisher, fields []interface{}) []interface{} {
	var address PublisherAddress
	p.Address = &address

	return append(
		fields, &address.Town, &address.Main, &address.Phone, &address.Email,
		&address.Site,
	)
}

var CollectionBuilders = map[string]func(*Publisher, []FilteringValue){
	"publisher_aliases": buildPublisherAliases,
}

var JoinBuilders = map[string]JoinBuilder{
	"publisher_addresses": buildPublisherAddress,
}
