package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
	
)

type Context struct{
	Writer http.ResponseWriter
	Req *http.Request
	Method string
	Path string
	Params map[string]string
	StatusCode int
	handlers []HandlerFunc
	index int
	engine *Engine
}

func NewContext(w http.ResponseWriter,r *http.Request) *Context{
	return &Context{
		Writer: w,
		Req: r,
		Method: r.Method,
		Path: r.URL.Path,
		index: -1,
		handlers: []HandlerFunc{},
	}
}

func (c *Context)Next(){
	c.index++
	
	if c.index<len(c.handlers){
		
		fmt.Printf("handle: handlers-index:%d\n", c.index)
		c.handlers[c.index](c)
	}else{
		fmt.Printf("index:%d,handler: NULL\n",c.index)
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

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	// 执行指定名称的模板并将结果写入响应
	// 参数说明:
	// - c.Writer: http.ResponseWriter 接口，用于输出模板渲染结果
	// - name: 要执行的模板名称 (对应模板定义中的 {{define "name"}})
	// - data: 传递给模板的数据对象 (在模板中通过 {{.}} 或 {{.FieldName}} 访问)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}