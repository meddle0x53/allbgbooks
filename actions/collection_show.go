package actions

import (
	"allbooks/models"
	"allbooks/routing"
)

func CollectionShowAction(context *routing.ResourceContext) {
	context.SetResponseData(models.CreateResource(context))
}
