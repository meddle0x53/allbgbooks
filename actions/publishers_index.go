package actions

import (
  "allbooks/models"
  "allbooks/routing"
)

func PublishersIndexAction(context *routing.CollectionContext) {
  publishers := models.GetPublishers(context.Page, context.PerPage, context.OrderBy)

  context.SetResponseData(publishers)
}
