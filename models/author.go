package models

type Author struct {
	ModelWithId
	Name        string `json:"name"`
	Nationality string `json:"nationality"`
}

func (author *Author) Fields() []interface{} {
	return []interface{}{&author.Id, &author.Name, &author.Nationality}
}

func GetAuthors(context CollectionContext) *[]Model {
	return CreateCollection(GetCollection("authors", context), "authors")
}
