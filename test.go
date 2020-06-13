package main

import (
	"fmt"
	"smarty/core"
)

func main() {
	var sm smarty.Smarty
	sm.Construct("tpl.config")
	sm.Assign("title", "首页")
	sm.Assign("menu", map[string]map[string]string{"index": map[string]string{"index": "菜单一", "url": "www.sogou.com"}, "blog": map[string]string{"index": "博客", "url": "www.sogou.com"}})
	sm.Assign("copyright", map[string]int{"time": 20190704})
	fmt.Println(sm.Display("index.htm"))
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
	}()
}
