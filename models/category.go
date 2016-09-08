package models

type Category struct {
	BaseModel
	Name string `json:"name"`
}

func (category *Category) Fields() []interface{} {
	return []interface{}{&category.Name}
}
