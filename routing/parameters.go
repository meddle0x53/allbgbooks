package routing

import (
  "strconv"
)

func ParseIntParam(ctx Context, name string, defaultVal uint64) uint64 {
  val, err := strconv.ParseUint(ctx.Request().Form.Get(name), 10, 64)
  if err != nil { val = defaultVal }

  return val
}
