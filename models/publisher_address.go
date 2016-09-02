package models

type PublisherAddress struct {
	Town  string `json:"town"`
	Main  string `json:"street"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Site  string `json:"site"`
}
