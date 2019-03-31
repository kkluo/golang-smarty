package main

import (
	"fmt"
	"smarty/core"
)

func main() {
	//*********添加读取配置文件tpl.config功能
	sm := smarty.Self{
		"tpl/",
		"{#",
		"$",
		"#}",
		true,
		"cache/",
		20,
	}
	sm.Assign("title", "首页")
	fmt.Println(sm.Display("index.htm"))
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
	}()
}
