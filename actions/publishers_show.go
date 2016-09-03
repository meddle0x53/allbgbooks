package actions

import (
	"allbooks/models"
	"allbooks/routing"
)

func PublishersShowAction(context *routing.ResourceContext) {
	publisher := models.GetPublisherById(context)

	context.SetResponseData(publisher)
	context.SetIsEmpty(publisher.Id == 0)
}
