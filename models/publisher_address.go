package models

type PublisherAddress struct {
	BaseModel
	Town  *string `json:"town"`
	Main  *string `json:"street"`
	Phone *string `json:"phone"`
	Email *string `json:"email"`
	Site  *string `json:"site"`
}

func (address *PublisherAddress) Fields() []interface{} {
	return []interface{}{
		&address.Town, &address.Main, &address.Phone, &address.Email, address.Site,
	}
}
