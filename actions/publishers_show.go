package actions

import (
	"allbooks/models"
	"allbooks/routing"
	"fmt"
)

func PublishersShowAction(context *routing.ResourceContext) {
	publisher := models.GetPublisherById(context.IdParameter)
	fmt.Println(publisher)

	context.SetResponseData(publisher)
	context.IsEmpty = publisher.Id == 0
}
