package actions

import (
	"allbooks/models"
	"allbooks/routing"
)

func PublishersShowAction(context *routing.ResourceContext) {
	publisher := models.GetPublisherById(context.IdParameter, context.JoinFields)

	context.SetResponseData(publisher)
	context.IsEmpty = publisher.Id == 0
}
