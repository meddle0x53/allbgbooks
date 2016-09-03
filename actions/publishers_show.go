package actions

import (
	"allbooks/models"
	"allbooks/routing"
)

func PublishersShowAction(context *routing.ResourceContext) {
	publisher := models.CreateResource(context, "publishers")

	context.SetResponseData(publisher)
}
