package smarty

//模板文件加载器,页面包含通过Load和include函数调用和迭代进行组装
import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

var readcache bool //是否读取了缓存文件
//加载模板文件
func (sm *Self) load(filepath string, re bool) string {
	readcache = false
	finalpath := ""
	//启用缓存判断缓存是否可用并加载缓存
	if sm.Caching == true {
		finalpath = sm.loadCache(filepath) //若存在缓存文件则返回缓存文件地址
	} else {
		finalpath = sm.Tpl_dir + filepath
	}
	fmt.Println(finalpath)
	//读取html文件
	html := sm.fileReader(finalpath)
	//未加载缓存的html导入编译器
	if readcache == false {
		html = sm.compile(html)
	}
	//fmt.Println(html)
	html = sm.include(html)
	//fmt.Println(filepath)
	//fmt.Println(html)
	//生成缓存文件
	if sm.Caching == true {
		cachepath := sm.Cache_dir + filepath + ".cache" //缓存路径
		//os.Remove(cachepath)
		err := ioutil.WriteFile(cachepath, []byte(html), 0666) //清空式写入
		if err != nil {
			panic(err)
		}
	}
	return html
}

//读取文件函数
func (sm *Self) fileReader(fpath string) string {
	html, err := ioutil.ReadFile(fpath)
	if err != nil {
		panic(err)
	}
	return string(html)
	//fmt.Println(html)
}

//判断缓存是否有效,有效返回缓存文件地址,无效则返回原tplpath地址
func (sm *Self) loadCache(tplpath string) string {
	//缓存文件情况
	tplpath = sm.Tpl_dir + tplpath
	cachepath := sm.Cache_dir + tplpath + ".cache"
	fstate, err := os.Stat(cachepath)
	//fmt.Println(fstate)
	if fstate == nil {
		readcache = false
		return tplpath
	}
	if err == nil {
		if sm.Cache_lifetime <= 0 { //无限期有效
			readcache = true
			return cachepath
		} else { //判断是否在有效期内
			/**fstate, err := f.Stat() //缓存文件时间
			if err != nil {
				panic("缓存文件信息错误")
			}**/
			overtime := (time.Now().Unix() - fstate.ModTime().Unix()) / 60
			//fmt.Println(fstate.ModTime().Unix())
			if overtime < sm.Cache_lifetime { //缓存文件生成时间超过时限
				readcache = false
				return tplpath
			} else {
				readcache = true
				return cachepath
			}
		}
	} else {
		return tplpath
	}
	/**f, err := os.Open(cachepath)
	defer f.Close()
	if err != nil { //打开缓存文件错误,缓存可能不存在
		return tplpath
	}**/

	return tplpath
}

//处理模板中的包含文件frame为框架页html字符串,返回替换了包含页（已编译）的html
func (sm *Self) include(frame string) string {
	tag := sm.Pre_tag + "include=(.*)" + sm.End_tag
	result := ""
	r, err := regexp.Compile(tag)
	if err == nil {
		allFile := r.FindAllStringSubmatch(frame, -1)
		for subtpl := range allFile {
			chtml := sm.load(allFile[subtpl][1], true) //已编译的subpage html
			fmt.Println(chtml)
			result = strings.Replace(frame, allFile[subtpl][0], chtml, -1) //替换掉原来页面的include语句
		}
	}
	return result

}
