package zweb

import (
	
	"strings"
)

type router struct {
	trees map[string]*node
}

type node struct {
	handler HandleFunc
	path    string
	children   map[string]*node
	starChild *node
	paramChild *node
	// route 记录到达该节点的完整路径
	route string
}

func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

func (r *router) addRoute(method, path string, handler HandleFunc) {
	tree, ok := r.trees[method]
	if !ok {
		tree = &node{path: "/"}
		r.trees[method] = tree
	}
	if path == "/" {
		tree.handler = handler
		return
	}
	cur := tree
	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		cur = cur.childOrCreate(seg)
	}
	cur.handler = handler
	cur.route = path
}

func (n *node) childOrCreate(path string) *node {
	if path == "*" {
		if n.starChild == nil {
			n.starChild = &node{
				path: path,
			}
		}
		return n.starChild
	}
	if path[0] == ':' {
		if n.paramChild == nil {
			n.paramChild = &node{
				path: path,
			}
		}
		return n.paramChild
	}
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]
	if !ok {
		child = &node{path: path}
		n.children[path] = child
	}
	return child
}

func (r *router) findRoute(method, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{n:root}, true
	}

	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")
	cur := root
	for _, seg := range segs {
		if cur.children == nil {
			if cur.paramChild != nil {
				mi := &matchInfo{
					n: cur.paramChild,
					pathParams: map[string]string{
						cur.paramChild.path[1:]: seg,
					},
				}
				return mi, true
			}
			return &matchInfo{n: cur.starChild}, cur.starChild != nil
		}
		child, ok := cur.children[seg]
		if !ok {
			if cur.paramChild != nil {
				mi := &matchInfo{
					n: cur.paramChild,
					pathParams: map[string]string{
						cur.paramChild.path[1:]: seg,
					},
				}
				return mi, true
			}
			return &matchInfo{n: cur.starChild}, cur.starChild != nil
		}
		cur = child
	}
	if cur.handler == nil {
		return nil, false
	}
	return &matchInfo{n: cur}, true
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}