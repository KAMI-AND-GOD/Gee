package gee
import(
	"time"
	"log"
)
func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s-%s in %v", c.StatusCode,c.Method, c.Req.RequestURI, time.Since(t))
	}
	
}