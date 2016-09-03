package models

type Author struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Nationality string `json:"nationality"`
}

func GetAuthors(context CollectionContext) *[]Model {
	return CreateCollection(GetCollection("authors", context), "authors")
}
