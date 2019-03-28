package smarty

//"fmt"
//"io/ioutil"
//"time"
//"strings"
//以loader为主线加载编译器并组装页面
//模板变量
var tplvals = make(map[string]interface{})

//smarty类,所有类名首字母大写
type Self struct {
	Tpl_dir        string //模板目录
	Pre_tag        string //模板前缀
	Var_tag        string //变量前缀如$
	End_tag        string //模板后缀
	Caching        bool   //是否启动缓存
	Cache_dir      string //缓存目录
	Cache_lifetime int64  //缓存有效时间（分钟）<=0时为无限期
}

func (sm *Self) Assign(key string, value interface{}) {
	_, exists := tplvals[key]
	if exists { //变量已经赋值
		panic("变量已经赋值")
	}
	tplvals[key] = value
}

//***显示模板即返回替换后的数组
//所有页面模块（include）模块分别进行缓存和预处理，最后再拼装
func (sm *Self) Display(tplname string) {
	sm.load(tplname, false) //
	//fmt.Println(page)
}
