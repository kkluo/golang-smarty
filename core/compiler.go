package smarty

import (
	//"encoding/json"
	//"reflect"
	"fmt"

	//"reflect"

	//"io/ioutil"
	"regexp"
	"strings"
)

type compiler struct { //编译类
}

var prochtml string //处理字符串

//编译处理器
//提取所有的标签再调用plugin里的处理模块
func (sm *Smarty) compile(html string) string {
	//fmt.Println(html)
	prochtml = html
	sm.easyReplace() //单变量替换
	//fmt.Println(prochtml)
	return prochtml
}

//单变量($menu.index.url)替换成字符串（根据系统设置进行转义）
//分为普通变量、数组变量
func (sm *Smarty) easyReplace() {
	rtag := "(?U)" + config.Pre_tag + "\\$(.*)" + config.End_tag //fmt.Println(rtag)
	r, err := regexp.Compile(rtag)
	if err == nil {
		valtags := r.FindAllStringSubmatch(prochtml, -1) //所有标签
		for _, valtag := range valtags {
			valkeys := strings.Split(valtag[1], ".") //键名数组
			if len(valkeys) < 1 {
				continue
			}
			tag := config.Pre_tag + "$" + valtag[1] + config.End_tag
			fmt.Println(tag)
			//mp := tplvals.(map[string]interface{})
			//fmt.Println(tplvals)
			val := fmt.Sprintf("%v", sm.getVals(tplvals, valkeys))
			fmt.Println(val)
			prochtml = strings.Replace(prochtml, tag, val, -1) //添加html处理
			//fmt.Println(val)
		}
	}
}

func (sm *Smarty) getVals(vals map[string]interface{}, keys []string) interface{} {
	for index, key := range keys {
		if val, exists := vals[key]; exists {
			if cmap, ok := val.(map[string]interface{}); ok {
				vals = cmap
			} else {
				fmt.Println(key, len(keys), index+1)
				if len(keys) == index+1 { //最后一次循环

					return val
				} else {
					//vals = val.(map[string]interface{})
					return nil
				}
			}
		} else {
			fmt.Println("键值", key)
			return nil
		}
	}
	return vals
}

//Interface to Map
func (sm *Smarty) interfaceToMap(from interface{}, to map[string]interface{}) {

}

//内容转义
func (sm *Smarty) stripHtml(html string) {
	//替换掉注释和一些标签
	reg := regexp.MustCompile(`<!--[^>]+>|<iframe[\S\s]+?</iframe>|<a[^>]+>|</a>|<script[\S\s]+?</script>|<div class="hzh_botleft">[\S\s]+?</div>`)
	html = reg.ReplaceAllString(html, "")
	//fmt.Println(html)
}
