package smarty

import (
	"fmt"
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

//单变量替换成字符串（根据系统设置进行转义）
//分为普通变量、数组变量
func (sm *Self) easyReplace() {
	tag := "(?U)" + sm.Pre_tag + "\\$(.*)" + sm.End_tag
	fmt.Println(tag)
	r, err := regexp.Compile(tag)
	if err == nil {
		valtags := r.FindAllStringSubmatch(prochtml, -1) //所有标签
		val := tplvals
		for key := range valtags {
			valkeys := strings.Split(valtags[key][1], ".") //键名数组
			fmt.Println(valkeys)

			for valkey := range valkeys {
				val = sm.getVals(val, string(valkey))
			}
			fmt.Println(val)
		}
	}
}
func (sm *Self) getVals(tplarr map[string]interface{}, key string) map[string]interface{} {
	if tplarr[key] != nil {
		return tplarr[key]
	} else {
		return ""
	}
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
