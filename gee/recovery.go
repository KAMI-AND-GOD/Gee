package gee
import(
	"runtime"
	"strings"
	"fmt"
	"log"
	"net/http"
)
// trace 打印当前goroutine的调用栈信息，用于调试
// message: 自定义的调试信息，会出现在调用栈前面
// 返回值: 包含调用栈信息的字符串
func trace(message string) string {
    // 创建一个足够大的数组来存储程序计数器（PC）
    var pcs [32]uintptr
    
    // runtime.Callers 获取当前goroutine的调用栈
    // 参数3表示跳过前3个调用者（包括runtime.Callers自身、trace函数和trace的调用者）
    n := runtime.Callers(3, pcs[:]) 

    var str strings.Builder
    str.WriteString(message + "\nTraceback:")
    
    // 遍历所有获取到的程序计数器
    for _, pc := range pcs[:n] {
        // 通过PC获取函数信息
        fn := runtime.FuncForPC(pc)
        // 获取该PC对应的文件名和行号
        file, line := fn.FileLine(pc)
        // 将信息格式化为 文件名:行号 的形式
        str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
    }
    return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}