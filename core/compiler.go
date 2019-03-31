package smarty

import (
	//"fmt"
	//"io/ioutil"
	"regexp"
	"strings"
)

var prochtml string //处理字符串

//编译处理器
//提取所有的标签再调用plugin里的处理模块
func (sm *Self) compile(html string) string {
	//fmt.Println(html)
	prochtml = html
	sm.easyReplace() //单变量替换
	//fmt.Println(prochtml)
	return prochtml
}

//单变量替换成字符串
func (sm *Self) easyReplace() {
	for key := range tplvals {
		tag := sm.Pre_tag + sm.Var_tag + key + sm.End_tag
		replace := tplvals[key].(string)                       //************************添加转义处理函数
		prochtml = strings.Replace(prochtml, tag, replace, -1) //替换变量
	}
	//fmt.Println(afthtml)
}

//循环处理
func (sm *Self) foreach() {

}

//内容转义
func (sm *Self) stripHtml(html string) {
	//替换掉注释和一些标签
	reg := regexp.MustCompile(`<!--[^>]+>|<iframe[\S\s]+?</iframe>|<a[^>]+>|</a>|<script[\S\s]+?</script>|<div class="hzh_botleft">[\S\s]+?</div>`)
	html = reg.ReplaceAllString(html, "")
	//fmt.Println(html)
}