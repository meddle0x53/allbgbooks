package models

type Genre struct {
	BaseModel
	Name string `json:"name"`
}

func (genre *Genre) Fields() []interface{} {
	return []interface{}{&genre.Name}
}
