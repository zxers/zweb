package zweb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func Test_router_addRoute(t *testing.T) {
	tests := []struct {
		name string
		// 输入
		method  string
		path    string
		handler HandleFunc
	}{
		{
			name:   "/",
			method: http.MethodGet,
			path:   "/",
		},
		{
			name:   "/",
			method: http.MethodGet,
			path:   "/user",
		},
		{
			name:   "/",
			method: http.MethodGet,
			path:   "//home",
		},
		{
			name:   "/",
			method: http.MethodGet,
			path:   "//home1///",
		},
		{
			name:   "/",
			method: http.MethodGet,
			path:   "/user/detail/profile",
		},
		{
			name:   "/",
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			name:   "/",
			method: http.MethodGet,
			path:   "/order/detail/:order_sn",
		},
	}
	var handler HandleFunc = func(ctx *Context) {

	}
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				handler: handler,
				path:    "/",
				children: map[string]*node{
					"user": {
						handler: handler,
						path:    "user",
						children: map[string]*node{
							"detail": {
								path: "detail",
								children: map[string]*node{
									"profile": {
										handler: handler,
										path:    "profile",
									},
								},
							},
						},
					},
					"home": {
						handler: handler,
						path:    "home",
					},
					"home1": {
						handler: handler,
						path:    "home1",
					},
					"order": {
						path:    "order",
						starChild: &node{
							path: "*",
							handler: handler,
						},
						children: map[string]*node{
							"detail": &node{
								path: "detail",
								paramChild: &node{
									path: ":order_sn",
									handler: handler,
								},
							},
						},
					},
				},
			},
		},
	}
	res := &router{map[string]*node{}}

	for _, tt := range tests {
		res.addRoute(tt.method, tt.path, handler)
	}
	errStr, ok := wantRouter.equal(res)
	assert.True(t, ok, errStr)

	findCases := []struct {
		name string
		// 输入
		method   string
		path     string
		handler  HandleFunc
		found    bool
		wantPath string
	}{
		{
			name:   "/",
			method: http.MethodGet,
			path:   "/",
			found:  true,
			wantPath: "/",
		},
		{
			name:   "/",
			method: http.MethodGet,
			path:   "/user",
			found:  true,
			wantPath: "user",
		},
		{
			name:   "/",
			method: http.MethodGet,
			path:   "/order/abs",
			found:  true,
			wantPath: "*",
		},
		{
			name:   "/",
			method: http.MethodGet,
			path:   "/order/detail/:232",
			found:  true,
			wantPath: ":order_sn",
		},
	}

	for _, tt := range findCases {
		node, ok := res.findRoute(tt.method, tt.path)
		assert.Equal(t, tt.found, ok)
		if !ok {
			return
		} 
		assert.Equal(t, tt.wantPath, node.n.path)
		assert.NotNil(t, node.n.handler)
	}
}

func (r router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		yv, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("目标 router 里面没有方法 %s 的路由树", k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return k + "-" + str, ok
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if y == nil {
		return "目标节点为 nil", false
	}
	if n.path != y.path {
		return fmt.Sprintf("%s 节点 path 不相等 x %s, y %s", n.path, n.path, y.path), false
	}

	nhv := reflect.ValueOf(n.handler)
	yhv := reflect.ValueOf(y.handler)
	if nhv != yhv {
		return fmt.Sprintf("%s 节点 handler 不相等 x %s, y %s", n.path, nhv.Type().String(), yhv.Type().String()), false
	}

	if len(n.children) != len(y.children) {
		return fmt.Sprintf("%s 子节点长度不等", n.path), false
	}
	if len(n.children) == 0 {
		return "", true
	}

	for k, v := range n.children {
		yv, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("%s 目标节点缺少子节点 %s", n.path, k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return n.path + "-" + str, ok
		}
	}
	return "", true
}
