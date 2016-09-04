package actions

import (
	"allbooks/routing"
)

func IndexAction(context routing.Context) {
	data := map[string]string{
		"publishersUrl": routing.Domain + "/publishers",
	}

	context.SetResponseData(data)
}
