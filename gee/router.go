package gee

import (
	"net/http"
	"strings"
)
func parsePath(path string) []string{
	parts:=make([]string,0)
	ss:=strings.Split(path,"/")
	for _,part:=range ss{
		if part!=""{
			parts=append(parts, part)
			if part[0]=='*'{
				break
			}
		}
	}
	return parts
}

type router struct {
	roots map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots: make(map[string]*node,0),
		handlers: make(map[string]HandlerFunc)}
}


func (r *router) addRoute(method string, path string, handler HandlerFunc) {
	key:=method + "-" + path
	r.handlers[key]=handler

	//前缀树中添加路由,目的是方便动态路由进行匹配
	parts:=parsePath(path)
	_,exsis:=r.roots[method]
	if !exsis{
		r.roots[method]=&node{part: method + "-"}
	}
	root:=r.roots[method]
	root.insert(path,parts,0)
}

//路由匹配: 从前缀树中找到与输入路由相匹配的第一条路由规则
func (r *router) getRoute(method string, path string) (string, map[string]string) {
	searchParts := parsePath(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return "", nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePath(n.path)
		for index, part := range parts {
			if part[0] == ':' {
				param:=part[1:]
				params[param] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				param:=part[1:]
				params[param] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n.path, params
	}

	return "", nil
}


func (r *router) handleReq(c *Context) {
	
	path,params:=r.getRoute(c.Method,c.Path)
	if path!=""{
		key:=c.Method +"-" + path
		c.Params=params
		handler:=r.handlers[key]
		c.handlers=append(c.handlers, handler)
	}else{
		c.handlers=append(c.handlers, func(c *Context){
			c.String(http.StatusNotFound, "404 NOT FOUND: %s-%s\n",c.Method, c.Path)
		})
	}
	//handle
	c.Next()
}