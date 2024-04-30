package zweb

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Context struct {
	Req    *http.Request
	Resp  http.ResponseWriter
	Params map[string]string
	MatchedRoute string

	RespStatusCode int
	RespData []byte

	tplEngine TemplateEngine
}

func (ctx *Context) BindJson(val interface{}) error {
	if ctx.Req.Body == nil {
		return errors.New("web:body为nil")
	}
	decode := json.NewDecoder(ctx.Req.Body)
	
	return decode.Decode(val)
}

func (ctx *Context) FormValue(key string) (string, error) {
	err := ctx.Req.ParseForm()
	if err != nil {
		return "", err
	}
	return ctx.Req.FormValue(key), nil
}

func (ctx *Context) QueryValue(key string) StringVal {
	params := ctx.Req.URL.Query()
	if params == nil {
		return StringVal{val: "", err: errors.New("web: 没有查询参数")}
	}
	vals, ok:= params[key]
	if !ok {
		return StringVal{val: "", err: errors.New("web: 找不到key")}
	}
	return StringVal{val: vals[0], err: nil}
}

func (ctx *Context) RespJSON(code int, val interface{}) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}
	ctx.RespStatusCode = code
	ctx.RespData = bs
	return err
}

func (ctx *Context) SetCookie(cookie http.Cookie) error {
	return nil
}

func (ctx *Context) Render(tplName string, data any) error {
	res, err := ctx.tplEngine.Render(ctx.Req.Context(), tplName, data)
	if err != nil {
		return err
	}
	ctx.RespData = res
	ctx.RespStatusCode = http.StatusOK
	return nil
}

func (ctx *Context) PathValue(key string) StringVal {
	val, ok := ctx.Params[key]
	if !ok {
		return StringVal{
			val: "",
			err: errors.New("key not found"),
		}
	}
	return StringVal{
		val: val,
		err: nil,
	}
}

type StringVal struct {
	val string
	err error
}

func (s StringVal) ToInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}

func (s StringVal) String() (string, error) {
	return s.val, s.err
}