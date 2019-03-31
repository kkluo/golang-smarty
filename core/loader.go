package smarty

//模板文件加载器,页面包含通过Load和include函数调用和迭代进行组装
import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	//"os"
	//"reflect"
	"regexp"
	"strings"
	//"time"
)

var readcache = make(map[string]bool)     //是否读取了缓存文件
var tplpath = make(map[string]string)     //模板文件地址
var cachepath = make(map[string]string)   //缓存文件地址
var nocache = make(map[string][][]string) //无缓存
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
	//在编译处理前提取nocache中内容
	if sm.Caching == true {
		sm.extractNocacheBlock(filepath, html)
	}
	//未加载缓存的html导入编译器
	if readcache[filepath] == false {
		html = sm.compile(html)
	}
	//fmt.Println(html)
	html = sm.include(html)
	//fmt.Println(html)
	//最后生成缓存文件
	if sm.Caching == true {
		sm.saveCache(filepath, html)
	}
	//return "" //测试
	//fmt.Println(html)
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

//处理模板中的包含文件frame为框架页html字符串,返回替换了包含页（已编译）的html
func (sm *Self) include(frame string) string {
	tag := "(?U)" + sm.Pre_tag + "include=(.*)" + sm.End_tag
	r, err := regexp.Compile(tag)
	if err == nil {
		allFile := r.FindAllStringSubmatch(frame, -1)
		//fmt.Println(allFile)
		for subfile := range allFile {
			//fmt.Println(allFile[subtpl][0])
			subpage := sm.load(allFile[subfile][1], true) //已编译的subpage html
			//fmt.Println(chtml)
			frame = strings.Replace(frame, allFile[subfile][0], subpage, -1) //替换掉原来页面的include语句
		}
	}
	return frame
}

//引擎最后一步去除页面所有smarty标签
func (sm *Self) stripSmartyTags(html string) string {
	tag := "(?U)" + sm.Pre_tag + "(.*)" + sm.End_tag
	r, err := regexp.Compile(tag)
	if err == nil {
		html = r.ReplaceAllString(html, "")
	}
	return html
}
