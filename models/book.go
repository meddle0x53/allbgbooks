package models

import "time"

type Book struct {
	ModelWithId
	ISBN        string     `json:"isbn"`
	Title       string     `json:"title"`
	Cover       string     `json:"cover"`
	Issue       int        `json:"issue"`
	Description *string    `json:"description"`
	Copies      *int       `json:"copies"`
	Price       *int       `json:"price"`
	PublishDate *time.Time `json:"publishDate"`
	Publisher   Model      `json:"publisher,omitempty"`
	Language    Model      `json:"language,omitempty"`
	Genre       Model      `json:"genre,omitempty"`
	Category    Model      `json:"category,omitempty"`
	Authors     *[]Model   `json:"authors,omitempty"`
}

func (book *Book) Fields() []interface{} {
	return []interface{}{
		&book.Id, &book.ISBN, &book.Title, &book.Cover, &book.Issue,
		&book.Description, &book.Copies, &book.Price, &book.PublishDate,
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
	case "categories":
		book.Category = (*collection)[0]
	case "authors":
		book.Authors = collection
	}
}
