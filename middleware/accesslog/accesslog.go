package accesslog

import (
	"encoding/json"
	"fmt"
	"zweb"
)

type MiddlewareBuilder struct {
	logFunc func(accessLog string)
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(accessLog string) {
			fmt.Println(accessLog)
		},
	}
}

func (m *MiddlewareBuilder) LogFunc(logFunc func(accessLog string)) *MiddlewareBuilder {
	m.logFunc = logFunc
	return m
}

func (m *MiddlewareBuilder) Build() zweb.Middleware {
	return func(next zweb.HandleFunc) zweb.HandleFunc {
		return func(ctx *zweb.Context) {
			defer func() {
				l := accessLog{
					Host:       ctx.Req.Host,
					Path:       ctx.Req.URL.Path,
					HTTPMethod: ctx.Req.Method,
				}
				val , _ := json.Marshal(l)
				m.logFunc(string(val))
			}()
			next(ctx)
		}
	}
}

type accessLog struct {
	Host       string
	Route      string
	HTTPMethod string `json:"http_method"`
	Path       string
}