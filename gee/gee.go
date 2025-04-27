package gee

import (
	"log"
	"net/http"
	"strings"
	"path"
	"html/template"
)
type HandlerFunc func(c *Context)

type H map[string]interface{}

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

	// htmlTemplates 是 HTML 模板集合，用于模板渲染
    // 使用标准库 template.Template 实现
    htmlTemplates *template.Template

    // funcMap 存储自定义模板函数，可在模板中使用
    funcMap  template.FuncMap
}

func New() *Engine{
	engine:=&Engine{router:*newRouter()}
	engine.RouteGroup=&RouteGroup{ engine: engine}
	engine.groups=append(engine.groups, engine.RouteGroup)
	return engine
}

// Default use Logger() & Recovery middlewares
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
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

func (group *RouteGroup) Use(middlewares ...HandlerFunc){
	group.middlewares=append(group.middlewares, middlewares...)
}

// createStaticHandler 创建静态文件处理函数
// relativePath: 路由前缀路径
// fs: 文件系统接口，用于访问静态文件
func (group *RouteGroup) createStaticHandler(RoutePath string, fs http.FileSystem) HandlerFunc {
    // 拼接完整的绝对路径（路由组前缀+相对路径）
    absolutePath := path.Join(group.prefix, RoutePath)
    
    // 创建文件服务器处理器，并去除请求路径中的前缀部分
    // http.StripPrefix 用于移除请求URL中的指定前缀
    // http.FileServer 是标准库提供的静态文件服务处理器
    fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
    
    // 返回自定义的HandlerFunc闭包
    return func(c *Context) {
        // 从URL参数中获取文件路径
        file := c.Param("filepath")
        
        // 检查文件是否存在及是否有访问权限
        if _, err := fs.Open(file); err != nil {
            c.Status(http.StatusNotFound)
            return
        }

        // 调用文件服务器处理HTTP请求
        fileServer.ServeHTTP(c.Writer, c.Req)
    }
}

// Static 注册静态文件路由
// relativePath: 路由前缀路径
// root: 静态文件在文件系统中的根目录
func (group *RouteGroup) Static(RoutePath string, root string) {
    // 创建静态文件处理函数
    // http.Dir 实现了http.FileSystem接口，用于访问本地文件系统
    handler := group.createStaticHandler(RoutePath, http.Dir(root))
    
    // 构建URL匹配模式，捕获文件路径参数
    // 例如: "/assets/*filepath" 可以匹配 "/assets/js/main.js"
    urlPattern := path.Join(RoutePath, "/*filepath")
    
    // 注册GET请求处理器
    group.GET(urlPattern, handler)
}

// SetFuncMap 设置自定义模板函数映射
// funcMap: 模板函数映射表，key 为函数名，value 为函数实现
// 注意：需要在加载模板前调用此方法
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
    engine.funcMap = funcMap
}

// LoadHTMLGlob 加载 HTML 模板文件
// pattern: 文件匹配模式，支持 glob 通配符
// 例如: "./templates/*.html"
// 该方法会解析模板并注册自定义函数，解析失败会 panic
func (engine *Engine) LoadHTMLGlob(pattern string) {
    // template.Must 包装模板解析，如果出错会 panic
    // template.New("") 创建新的模板集合
    // Funcs() 注册自定义模板函数
    // ParseGlob() 根据通配符模式加载模板文件
    engine.htmlTemplates = template.Must(
        template.New("").Funcs(engine.funcMap).ParseGlob(pattern),
    )
}

func (e *Engine) Run(port string){
	http.ListenAndServe(port,e)//e.ServeHTTP()会在http请求过来后自动调用
}

//实现http.ListenAndServe(string , http.Handler)
//http.Handler接口 需要实现ServeHTTP(w http.ResponseWriter,r *http.Request)
func (e *Engine) ServeHTTP(w http.ResponseWriter,r *http.Request){
	var middlewares []HandlerFunc
	for _,group:=range e.groups{
		if strings.HasPrefix(r.URL.Path,group.prefix){
			middlewares=append(middlewares, group.middlewares...)
		}
	}
	c:=NewContext(w,r)
	c.engine = e
	c.handlers=middlewares
	e.router.handleReq(c)
}
