package main
import(
	"gee"
	"net/http"
)
func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	r.GET("/hello", func(c *gee.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s! Your age is %s, you're at %s\n", c.Query("name"),c.Query("age"),c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.ParseForm("username"),
			"password": c.ParseForm("password"),
		})
	})
	r.Run(":9999")
}