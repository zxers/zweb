package cookie

import (
	"net/http"
	"zweb/session"
)

var _ session.Propagator = &Propagator{}

type Propagator struct {
	cookieName string
}

func NewPropagator() *Propagator {
	return &Propagator{
		cookieName: "_sessid",
	}
}

// Extract implements session.Propagator.
func (p *Propagator) Extract(req *http.Request) (string, error) {
	cookie, err := req.Cookie(p.cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// Inject implements session.Propagator.
func (p *Propagator) Inject(id string, writer http.ResponseWriter) error {
	http.SetCookie(writer, &http.Cookie{
		Name: p.cookieName,
		Value: id,
	})
	return nil
}

// Remove implements session.Propagator.
func (p *Propagator) Remove(writer http.ResponseWriter) error {
	http.SetCookie(writer, &http.Cookie{
		Name: p.cookieName,
		MaxAge: -1,
	})
	return nil
}
