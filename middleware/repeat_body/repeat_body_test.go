package repeatbody

import (
	"testing"
	"zweb"
)

func TestMiddleware(t *testing.T) {
	s := zweb.NewHTTPServer()
	s.Use(Middleware())
	s.Get("/user", func(ctx *zweb.Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Start(":8084")
}