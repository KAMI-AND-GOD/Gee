package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)
type H map[string]interface{}

type Context struct{
	Writer http.ResponseWriter
	Req *http.Request
	Method string
	Path string
	Params map[string]string
	StatusCode int
}

func NewContext(w http.ResponseWriter,r *http.Request) *Context{
	return &Context{
		Writer: w,
		Req: r,
		Method: r.Method,
		Path: r.URL.Path,
	}
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) Param(key string)string{
	return c.Params[key]
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) Query(key string) string{
	return c.Req.URL.Query().Get(key)
}

func (c *Context) ParseForm(key string) string{
	return c.Req.FormValue(key)
}

func (c *Context) Data(code int,data string){
	c.Status(code)
	c.Writer.Write([]byte(data))
}

func (c *Context) String(code int ,format string,values...interface{}){
	c.Status(code)
	c.SetHeader("Content-Type", "text/plain")
	c.Writer.Write([]byte(fmt.Sprintf(format,values...)))
}

func (c *Context) JSON(code int, js interface{}){
	c.Status(code)
	c.SetHeader("Content-Type","application/json")
	encoder:=json.NewEncoder(c.Writer)
	if err:=encoder.Encode(js);err!=nil{
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
