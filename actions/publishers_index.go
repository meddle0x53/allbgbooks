package actions

import (
	"allbooks/models"
	"allbooks/routing"
)

func PublishersIndexAction(context *routing.CollectionContext) {
	context.SetResponseData(models.GetPublishers(context))
}
