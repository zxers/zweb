package session

import (
	"fmt"
	"net/http"
	"zweb"
	"testing"
)

func TestType(t *testing.T) {
	s := zweb.NewHTTPServer()
	m := &Manager{

	}
	s.Use(func(next zweb.HandleFunc) zweb.HandleFunc {
		return func(ctx *zweb.Context) {
			if ctx.Req.URL.Path != "/login" {
				session, err := m.GetSession(ctx)
				if err != nil {
					ctx.RespStatusCode = http.StatusUnauthorized
					return
				}
				fmt.Println(session)
				m.RefreshSession(ctx)
			}
		}
	})

	s.Get("/login", func(ctx *zweb.Context) {
		// 登录成功后
		session, err := m.InitSession(ctx)
		if err != nil {
			return
		}
		session.Set(ctx.Req.Context(), "key", "val")
	})

	s.Get("/resource", func(ctx *zweb.Context) {
		// 登录成功后
		session, err := m.GetSession(ctx)
		if err != nil {
			return
		}
		
		val, err := session.Get(ctx.Req.Context(), "key")
		if err != nil {
			return
		}
		v, ok := val.([]byte)
		if !ok {
			return
		}
		ctx.RespData = v
	})

	s.Get("/logout", func(ctx *zweb.Context) {
		// 登录成功后
		err := m.RemoveSession(ctx)
		if err != nil {
			return
		}
	})

	s.Start(":8082")
}