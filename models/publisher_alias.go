package models

type PublisherAlias struct {
	BaseModel
	Name string `json:"name"`
}

func (publisherAlias *PublisherAlias) Fields() []interface{} {
	return []interface{}{&publisherAlias.Name}
}
