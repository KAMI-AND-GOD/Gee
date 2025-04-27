package main
import(
	"gee"
	"net/http"
)
func main(){
	r:=gee.New()
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})
	r.GET("/hello", func(c *gee.Context) {
		// expect /hello?name=kami
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	
	r.Run(":9999")
}