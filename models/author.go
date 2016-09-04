package models

import (
	"strconv"
)

type Author struct {
	BaseModel
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Nationality string `json:"nationality"`
}

func (author *Author) Fields() []interface{} {
	return []interface{}{&author.Id, &author.Name, &author.Nationality}
}

func (author *Author) Identifier() string {
	return strconv.Itoa(author.Id)
}

func (author *Author) IsEmpty() bool {
	return author.Id <= 0
}

func GetAuthors(context CollectionContext) *[]Model {
	return CreateCollection(GetCollection("authors", context), "authors")
}
