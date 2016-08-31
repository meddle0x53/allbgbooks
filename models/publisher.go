package models

type Publisher struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Code  string `json:"code"`
	State string `json:"state"`
}

type Publishers []Publisher

func GetPublishers(
	page uint64, perPage uint64, orderBy string,
	filters []FilteringValue, ignoreCase bool,
) *Publishers {
	rows := GetCollection(
		"publishers", page, perPage, orderBy, filters, ignoreCase,
	)
	defer rows.Close()

	result := Publishers{}
	for rows.Next() {
		var publisher Publisher

		err := rows.Scan(
			&publisher.Id, &publisher.Name, &publisher.Code, &publisher.State)
		if err != nil {
			panic(err)
		}

		result = append(result, publisher)
	}

	return &result
}

func GetPublisherById(id string) *Publisher {
	var publisher Publisher

	GetResource("publishers", id).Scan(
		&publisher.Id, &publisher.Name, &publisher.Code, &publisher.State)

	return &publisher
}
