package main

import (
  "allbooks/actions"
)

func AppRoutes () Routes {
  return Routes{
    Route{
      "Index",
      "GET",
      "/",
      actions.IndexAction,
    },
    Route{
      "Publishers",
      "GET",
      "/publishers",
      actions.PublishersIndexAction,
    },
  }
}
