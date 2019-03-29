package auditlog

import (
	"net/http"
)

type Audit struct{}

func NewAudit() *Audit {
	return &Audit{}
}

func (a *Audit) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		// todo: before business logic
		next(rw, r)
		// todo: after business logic
	default:
		next(rw, r)
	}
}
