package models

type Language struct {
	BaseModel
	Name string `json:"name"`
}

func (language *Language) Fields() []interface{} {
	return []interface{}{&language.Name}
}
