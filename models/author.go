package models

type Author struct {
	ModelWithId
	Name        string `json:"name"`
	Nationality string `json:"nationality"`
}

func (author *Author) Fields() []interface{} {
	return []interface{}{&author.Id, &author.Name, &author.Nationality}
}
