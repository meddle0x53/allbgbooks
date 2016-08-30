package actions

import (
  "allbooks/routing"
)

func PublishersShowAction(context *routing.ResourceContext) {
  context.SetResponseData(make(map[string]string))
}
