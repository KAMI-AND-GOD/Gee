package gee

import(
	"fmt"
	"net/http"
)

type Engine struct{
	router map[string]http.HandlerFunc
}

//实现http.ListenAndServe(string , http.Handler)
//http.Handler接口 需要实现ServeHTTP(w http.ResponseWriter,r *http.Request)
func (e *Engine) ServeHTTP(w http.ResponseWriter,r *http.Request){
	method:=r.Method
	uri:=r.URL.Path
	key:=method +"-"+ uri
	handlerFunc,exsis:=e.router[key]
	if exsis{
		handlerFunc(w,r)
	}else {
		w.WriteHeader(404)
		fmt.Fprintf(w,"Server %s: 404 not found!\n",key)
	}
}

func New() *Engine{
	return &Engine{
		router: make(map[string]http.HandlerFunc)}
}

func (e *Engine) GET(uri string,hf http.HandlerFunc){
	key:="GET"+"-"+ uri
	e.router[key]=hf
}

func (e *Engine) POST(uri string,hf http.HandlerFunc){
	key:="POST"+"-"+ uri
	e.router[key]=hf
}

func (e *Engine) Run(port string){
	http.ListenAndServe(port,e)//e.ServeHTTP()会在http请求过来后自动调用
}

