package smarty

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//"time"
//"strings"
//以loader为主线加载编译器并组装页面
//type vals interface{}
var tplvals = make(map[string]interface{}) //模板变量

//为了提升加载速度当包含页未修改被包含页修改后，则不会自动更新包含页缓存，只有包含页更新后才会重新检查被包含页情况
//程序开发过程中必须将Caching设为false即取消缓存
type Smarty struct { //smarty类,模板设置,所有类名首字母大写
	Tpl_dir        string `json:"Tpl_dir"`        //模板目录
	Pre_tag        string `json:"Pre_tag"`        //模板前缀
	End_tag        string `json:"End_tag"`        //模板后缀
	Caching        bool   `json:"Caching"`        //是否启动缓存
	Cache_dir      string `json:"Cache_dir"`      //缓存目录
	Cache_lifetime int64  `json:"Cache_lifetime"` //缓存有效时间（分钟）<=0时为无限期
}

var config Smarty

//模板配置文件tpl.config功能
func (sm Smarty) Construct(filename string) {
	//var config = new(Smarty)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("配置文件打开错误"+filename, err)
		panic(1)
	}
	err = json.Unmarshal([]byte(data), &config)
	if err != nil {
		fmt.Println("配置文件格式错误"+filename, err)
	}
	fmt.Println(config)
}
func (sm Smarty) Assign(key string, value interface{}) {
	_, exists := tplvals[key]
	if exists { //变量已经赋值
		panic("变量已经赋值")
	}
	tplvals[key] = value
}

//***显示模板即返回替换后的数组
//所有页面模块（include）模块分别进行缓存和预处理，最后再拼装
func (sm Smarty) Display(tplname string) string {
	html := sm.load(tplname, false) //
	html = sm.stripSmartyTags(html)
	return html
	//fmt.Println(page)
}
