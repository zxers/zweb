package zweb

import (
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	s.Get("/", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, user"))
	})
	s.Get("/user/*", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, user ****"))
	})
	s.Get("/user/home/:id", func(ctx *Context) {
		ctx.Resp.Write([]byte(fmt.Sprintf("hhh, %s", ctx.Params["id"])))
	})
	s.Start(":8082")
}