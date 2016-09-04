package actions

import (
	"allbooks/models"
	"allbooks/routing"
)

func AuthorsShowAction(context *routing.ResourceContext) {
	author := models.CreateResource(context, "authors")

	context.SetResponseData(author)
}
