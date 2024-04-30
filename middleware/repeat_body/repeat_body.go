package repeatbody

import (
	"io"
	"zweb"
)

func Middleware() zweb.Middleware {
	return func(next zweb.HandleFunc) zweb.HandleFunc {
		return func(ctx *zweb.Context) {
			ctx.Req.Body = io.NopCloser(ctx.Req.Body)
			ctx.Resp.Write([]byte("start"))
			next(ctx)
			ctx.Resp.Write([]byte("end"))
		}
	}
}