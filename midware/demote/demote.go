package demote

import "net/http"

type Demote struct{}

func NewDemote() *Demote {
	return &Demote{}
}

func (d *Demote) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// todo: demote
	next(rw, r)
}
