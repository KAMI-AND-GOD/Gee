package gee
import(
	"strings"
)
type node struct{
	part string
	path string
	isWild bool
	children []*node
}
func (n *node) matchChild(part string) *node{
	for _,child:=range n.children{
		if child.part==part||child.isWild{
			return child
		}
	}
	return nil

}
func (n *node) matchChildren(part string) []*node{
	children:=[]*node{}
	for _,child:=range n.children{
		if child.part==part||child.isWild{
			children=append(children, child)
		}
	}
	return children
}
func (n *node) insert(path string, parts []string, height int){//height从0开始
	if height>len(parts)-1{
		n.path=path //叶子结点存完整路由：如/hello/:name
		return
	}
	part:=parts[height]
	child:=n.matchChild(part)
	if child==nil{
		child=&node{
			part: part,
			isWild: part[0]==':'||part[0]=='*',
		}
		n.children=append(n.children, child)
	}
	child.insert(path,parts,height+1)
}

//返回存有匹配路由的叶子结点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.path == "" {
			return nil //跟已有的路由规则匹配了前缀，但并不完整，所有不是合法的路由
						//举个例子：设置了一个路由规则为/hello/:name/:age ，输入路由为/hello/kami，只匹配到/hello/:name，需要/hello/kami/18
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}