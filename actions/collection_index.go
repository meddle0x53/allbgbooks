package actions

import (
	"allbooks/models"
	"allbooks/routing"
)

func CollectionIndexAction(context *routing.CollectionContext) {
	name := context.CollectionName()
	collection := models.CreateCollection(models.GetCollection(context), name)
	context.SetResponseData(collection)
}
