package actions

import (
	"allbooks/models"
	"allbooks/routing"
)

func AuthorsIndexAction(context *routing.CollectionContext) {
	context.SetResponseData(models.GetAuthors(context))
}
