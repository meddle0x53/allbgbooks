package models

type PublisherContact struct {
	BaseModel
	Name string `json:"name"`
}

func (publisherContact *PublisherContact) Fields() []interface{} {
	return []interface{}{&publisherContact.Name}
}
