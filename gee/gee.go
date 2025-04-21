package gee

import(
	"net/http"
)
type HandlerFunc func(c *Context)

type Engine struct{
	router router
}

func New() *Engine{
	return &Engine{
		router: *newRouter(),
	}
}

//实现http.ListenAndServe(string , http.Handler)
//http.Handler接口 需要实现ServeHTTP(w http.ResponseWriter,r *http.Request)
func (e *Engine) ServeHTTP(w http.ResponseWriter,r *http.Request){
	c:=NewContext(w,r)
	e.router.handleReq(c)
}

func (e *Engine) GET(path string,hf HandlerFunc){
	key:="GET"+"-"+ path
	e.router.handlers[key]=hf
}

func (e *Engine) POST(path string,hf HandlerFunc){
	key:="POST"+"-"+ path
	e.router.handlers[key]=hf
}

func (e *Engine) Run(port string){
	http.ListenAndServe(port,e)//e.ServeHTTP()会在http请求过来后自动调用
}

