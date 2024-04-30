package zweb

import (
	"net"
	"net/http"
)

type HandleFunc func(*Context)

type Server interface {
	http.Handler
	Start(addr string) error
	Addroute(method, path string, handler HandleFunc)
}

type HTTPServer struct {
	router
	mdls []Middleware
	tplEngine TemplateEngine
}

type ServerOption func(s *HTTPServer)

func ServerWithTemplateEngine(tplEngine TemplateEngine) ServerOption {
	return func(s *HTTPServer) {
		s.tplEngine = tplEngine
	}
}

func NewHTTPServer(opts ...ServerOption) *HTTPServer {
	hs :=  &HTTPServer{
		router: newRouter(),
	}

	for _, opt := range opts {
		opt(hs)
	}

	return hs
}

func (s *HTTPServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Req: req,
		Resp: resp,
		tplEngine: s.tplEngine,
	}
	// 查找路由
	// 执行handler
	root := s.serve
	for i := len(s.mdls) - 1; i >= 0; i-- {
		m := s.mdls[i]
		root = m(root)
	}
	var flash Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			s.flash(ctx)
		}
	}
	root = flash(root)
	root(ctx)

}

func (s *HTTPServer) serve(ctx *Context) {
	mi, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	
	if !ok || mi.n.handler == nil {
		ctx.RespStatusCode = http.StatusNotFound 
		ctx.RespData = []byte("404没找到")
		return
	}

	ctx.Params = mi.pathParams
	ctx.MatchedRoute = mi.n.route
	mi.n.handler(ctx)
}

func (s *HTTPServer) Start(addr string) error {

	l, _ := net.Listen("tcp", addr)

	return http.Serve(l, s)
}

func (s *HTTPServer) Addroute(method, path string, handler HandleFunc) {
	s.addRoute(method, path, handler)
}

func (s *HTTPServer) Get(path string, handler HandleFunc) {
	s.Addroute("GET", path, handler)
}

func (s *HTTPServer) Use(mdls ...Middleware) {
	s.mdls = append(s.mdls, mdls...)
}

func (s *HTTPServer) flash(ctx *Context) {
	ctx.Resp.WriteHeader(ctx.RespStatusCode)
	ctx.Resp.Write(ctx.RespData)
}

type HTTPSServer struct {
	Server
	CertFile string
	KeyFile string
}

func (s *HTTPSServer) Start(addr string) error {
	return http.ListenAndServeTLS(addr, s.CertFile, s.KeyFile, s)
}

