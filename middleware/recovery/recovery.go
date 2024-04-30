package recovery

import "zweb"

type MiddlewareBuilder struct {
	StatusCode int
	ErrMsg string
	LogFunc func(ctx *zweb.Context)
}

func NewMiddlewareBuilder(statusCode int, errMsg string, logFunc func(ctx *zweb.Context)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		StatusCode: statusCode,
		ErrMsg: errMsg,
		LogFunc: logFunc,
	}
}

func (m *MiddlewareBuilder) Build() zweb.Middleware {
	return func(next zweb.HandleFunc) zweb.HandleFunc {
		return func(ctx *zweb.Context) {
			defer func(){
				if err := recover(); err != nil {
					ctx.RespStatusCode = m.StatusCode
					ctx.RespData = []byte(m.ErrMsg)
					m.LogFunc(ctx)
				}
			}()
			next(ctx)
		}
	}
}