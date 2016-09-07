package models

type Book struct {
	ModelWithId
	ISBN        string  `json:"isbn"`
	Title       string  `json:"title"`
	Cover       string  `json:"cover"`
	Issue       int     `json:"issue"`
	Description *string `json:"description"`
	Copies      *int    `json:"copies"`
	Publisher   Model   `json:"publisher,omitempty"`
	Language    Model   `json:"language,omitempty"`
	Genre       Model   `json:"genre,omitempty"`
}

func (book *Book) Fields() []interface{} {
	return []interface{}{
		&book.Id, &book.ISBN, &book.Title, &book.Cover, &book.Issue,
		&book.Description, &book.Copies,
	}
}

func (book *Book) SetRelation(name string, collection *[]Model) {
	switch name {
	case "publishers":
		book.Publisher = (*collection)[0]
	case "languages":
		book.Language = (*collection)[0]
	case "genres":
		book.Genre = (*collection)[0]
	}
}
