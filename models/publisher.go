package models

type Publisher struct {
	ModelWithId
	Name     string   `json:"name"`
	Code     string   `json:"code"`
	State    string   `json:"state"`
	Address  Model    `json:"address,omitempty"`
	Aliases  *[]Model `json:"aliases,omitempty"`
	Contacts *[]Model `json:"contacts,omitempty"`
}

func (publisher *Publisher) Fields() []interface{} {
	return []interface{}{
		&publisher.Id, &publisher.Name, &publisher.Code, &publisher.State,
	}
}

func (publisher *Publisher) SetRelation(name string, collection *[]Model) {
	switch name {
	case "publisher_addresses":
		publisher.Address = (*collection)[0]
	case "publisher_aliases":
		publisher.Aliases = collection
	case "publisher_contacts":
		publisher.Contacts = collection
	}
}

func GetPublishers(context CollectionContext) *[]Model {
	return CreateCollection(GetCollection("publishers", context), "publishers")
}
