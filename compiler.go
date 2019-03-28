package smarty

import (
	//"fmt"
	//"io/ioutil"
	"strings"
)

var prochtml string //处理字符串

//编译处理器
func (sm *Self) compile(html string) string {
	prochtml = html
	sm.easyReplace()
	//fmt.Println(prohtml)
	return prochtml
}

//将模板中所有的变量替换成数据
func (sm *Self) easyReplace() {
	for key := range tplvals {
		tag := sm.Pre_tag + sm.Var_tag + key + sm.End_tag
		prochtml = strings.Replace(prochtml, tag, tplvals[key].(string), -1) //替换变量
	}
	//fmt.Println(afthtml)
}
