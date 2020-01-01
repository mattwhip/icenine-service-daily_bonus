package grifts

import (
	"github.com/mattwhip/icenine-service-daily_bonus/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
