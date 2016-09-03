package routing

import (
	"strconv"
)

func ParseIntParam(ctx Context, name string, defaultVal uint64) uint64 {
	val, err := strconv.ParseUint(ctx.Request().Form.Get(name), 10, 64)
	if err != nil {
		val = defaultVal
	}

	return val
}

func ParseSimpleIntParam(ctx Context, name string, defaultVal int) int {
	val, err := strconv.Atoi(ctx.Request().Form.Get(name))
	if err != nil {
		val = defaultVal
	}

	return val
}
