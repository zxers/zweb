package errhdl

import "zweb"

type MiddlewareBuilder struct {
	resp map[int][]byte
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		resp: map[int][]byte{},
	}
}

func (m *MiddlewareBuilder) RegisterErr(code int, resp []byte) *MiddlewareBuilder {
	m.resp[code] = resp
	return m
}

func (m *MiddlewareBuilder) Build() zweb.Middleware {
	return func(next zweb.HandleFunc) zweb.HandleFunc {
		return func(ctx *zweb.Context) {
			next(ctx)
			resp, ok := m.resp[ctx.RespStatusCode]
			if ok {
				ctx.RespData = resp
			}
		}
	}
}
