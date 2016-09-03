package models

type PublisherAlias struct {
	Name string `json:"name"`
}

func (publisherAlias *PublisherAlias) Fields() []interface{} {
	return []interface{}{&publisherAlias.Name}
}
