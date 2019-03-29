package midware

import (
	"github.com/chinx/coupon/midware/auditlog"
	"github.com/chinx/coupon/midware/demote"
	"github.com/chinx/coupon/midware/ratelimit"
	"github.com/chinx/coupon/midware/session"
	"github.com/urfave/negroni"
)

func SetMidwares(n *negroni.Negroni) {
	n.Use(auditlog.NewAudit())
	n.Use(ratelimit.NewRateLimit())
	n.Use(demote.NewDemote())
	n.Use(session.NewSession("/v1/user/login"))
}
