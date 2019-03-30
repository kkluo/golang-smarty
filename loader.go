package smarty

//模板文件加载器,页面包含通过Load和include函数调用和迭代进行组装
import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

var readcache = make(map[string]bool)   //是否读取了缓存文件
var tplpath = make(map[string]string)   //模板文件地址
var cachepath = make(map[string]string) //缓存文件地址
//加载模板文件
func (sm *Self) load(filepath string, re bool) string {
	tplpath[filepath] = sm.Tpl_dir + filepath
	cachepath[filepath] = sm.Cache_dir + base64.URLEncoding.EncodeToString([]byte(filepath)) + ".cache" //模板文件采用urlencode转义
	if sm.Caching == true {                                                                             //只有启动缓存时才检查缓存文件
		readcache[filepath] = sm.isCached(filepath)
	} else {
		readcache[filepath] = false
	}
	finalpath := ""
	//fmt.Println(filepath)
	//启用缓存判断缓存是否可用并加载缓存
	if readcache[filepath] {
		finalpath = cachepath[filepath] //若存在缓存文件则返回缓存文件地址
	} else {
		finalpath = tplpath[filepath]
	}
	//fmt.Println(finalpath)
	html := sm.fileReader(finalpath) //读取文件
	//fmt.Println(readcache[filepath])
	//未加载缓存的html导入编译器
	if readcache[filepath] == false {
		html = sm.compile(html)
	}
	//fmt.Println(html)
	html = sm.include(html)
	//fmt.Println(html)
	//生成缓存文件
	if sm.Caching == true && readcache[filepath] == false {
		//fmt.Println("更新了缓存" + filepath)
		err := ioutil.WriteFile(cachepath[filepath], []byte(html), 0666) //清空式写入
		if err != nil {
			panic("缓存写入失败" + fmt.Sprintf("%s", err))
		}
	}
	return html
}

//读取文件函数
func (sm *Self) fileReader(fpath string) string {
	html, err := ioutil.ReadFile(fpath)
	if err != nil {
		panic("模板或缓存文件读取失败" + fmt.Sprintf("%s", err))
	}
	return string(html)
	//fmt.Println(html)
}

//判断缓存是否有效,有效返回cache文件地址,无效则返回原tplpath地址
func (sm *Self) isCached(srcpath string) bool {
	//fmt.Println("缓存地址" + srcpath)
	//判断缓存生成时间是否超限
	cstate, _ := os.Stat(cachepath[srcpath]) //缓存文件对象状态
	if cstate == nil {                       //缓存文件不存在
		return false
	}
	tstate, _ := os.Stat(tplpath[srcpath]) //模板文件对象状态
	if tstate == nil {                     //模板文件不存在
		panic("模板文件不存在")
	}
	if cstate.ModTime().Unix() < tstate.ModTime().Unix() { //模板发生修改后不读取缓存
		//fmt.Println("模板文件已经发生修改")
		return false
	}
	if sm.Cache_lifetime <= 0 { //无限期有效
		return true
	} else { //判断是否在有效期内
		overtime := (time.Now().Unix() - cstate.ModTime().Unix()) / 60
		if overtime > sm.Cache_lifetime { //缓存文件生成时间超过时限
			//fmt.Println(overtime)
			return false
		} else {
			return true
		}
	}
}

//处理模板中的包含文件frame为框架页html字符串,返回替换了包含页（已编译）的html
func (sm *Self) include(frame string) string {
	tag := sm.Pre_tag + "include=(.*)" + sm.End_tag
	result := frame
	r, err := regexp.Compile(tag)
	if err == nil {
		allFile := r.FindAllStringSubmatch(frame, -1)
		//fmt.Println(allFile)
		for subtpl := range allFile {
			//fmt.Println(allFile[subtpl][0])
			chtml := sm.load(allFile[subtpl][1], true) //已编译的subpage html
			//fmt.Println(chtml)
			result = strings.Replace(frame, allFile[subtpl][0], chtml, -1) //替换掉原来页面的include语句
			frame = result
		}
	}
	return result
}
