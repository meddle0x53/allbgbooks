package actions

import (
  "net/http"
  "github.com/gorilla/context"
  "database/sql"
  "allbooks/models"
  "allbooks/decorators"
)

func PublishersIndexAction(w http.ResponseWriter, r *http.Request) {
  db := context.Get(r, "dbConnection").(*sql.DB)
  publishers := models.GetPublishers(db)

  params := context.Get(r, "params").(*decorators.Params)
  RenderJson(publishers, w, params.Pretty)
}
