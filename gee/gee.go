package gee

import(
	"net/http"
	"log"
)
type HandlerFunc func(c *Context)

type RouteGroup struct{
	prefix string
	parent *RouteGroup
	middlewares []HandlerFunc
	engine *Engine
	
}

type Engine struct{
	*RouteGroup
	router router
	groups []*RouteGroup
}

func New() *Engine{
	engine:=&Engine{router:*newRouter()}
	engine.RouteGroup=&RouteGroup{ engine: engine}
	engine.groups=append(engine.groups, engine.RouteGroup)
	return engine
}

func (group *RouteGroup) Group(prefix string)*RouteGroup{
	newGroup:=&RouteGroup{
		parent: group,
		prefix: group.prefix + prefix,
		engine: group.engine,
	}
	group.engine.groups=append(group.engine.groups, newGroup)
	return newGroup
}

func (group *RouteGroup) GET(path string,hf HandlerFunc){
	newPath:=group.prefix + path
	log.Printf("Route %4s - %s", "GET", newPath)
	group.engine.router.addRoute("GET",newPath,hf)
}

func (group *RouteGroup) POST(path string,hf HandlerFunc){
	newPath:=group.prefix + path
	log.Printf("Route %4s - %s", "POST", newPath)
	group.engine.router.addRoute("POST",newPath,hf)
}

func (e *Engine) Run(port string){
	http.ListenAndServe(port,e)//e.ServeHTTP()会在http请求过来后自动调用
}

//实现http.ListenAndServe(string , http.Handler)
//http.Handler接口 需要实现ServeHTTP(w http.ResponseWriter,r *http.Request)
func (e *Engine) ServeHTTP(w http.ResponseWriter,r *http.Request){
	c:=NewContext(w,r)
	e.router.handleReq(c)
}
