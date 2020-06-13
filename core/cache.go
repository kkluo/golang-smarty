package smarty

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

//判断缓存是否有效,有效返回cache文件地址,无效则返回原tplpath地址
func (sm *Smarty) isCached(srcpath string) bool {
	//return false //缓存生成测试
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
	if config.Cache_lifetime <= 0 { //无限期有效
		return true
	} else { //判断是否在有效期内
		overtime := (time.Now().Unix() - cstate.ModTime().Unix()) / 60
		if overtime > config.Cache_lifetime { //缓存文件生成时间超过时限
			//fmt.Println(overtime)
			return false
		} else {
			return true
		}
	}
}

//从页面中提取nocache内容并保存到字典，在保存缓存时替换掉nocache中内容
func (sm *Smarty) extractNocacheBlock(tplpath string, html string) {
	preg := "(?U)nocache" + config.End_tag + "([\\s\\S]*)" + config.Pre_tag + "/nocache"
	r, err := regexp.Compile(preg)
	if err == nil {
		nocache[tplpath] = r.FindAllStringSubmatch(html, -1) //所有不缓存内容
	}
	//fmt.Println(len(nocache[tplpath]))
	//fmt.Println(nocache[tplpath])
}

//将nocache中已编译内容替换回模板中内容，即调用缓存的时候仍需要编译实现nocache效果
func (sm *Smarty) nocacheRollback(filepath string, html string) string {
	stag := config.Pre_tag + "nocache" + config.End_tag
	etag := config.Pre_tag + "/nocache" + config.End_tag
	for i := range nocache[filepath] {
		html = strings.Replace(html, html[strings.Index(html, stag):strings.Index(html, etag)+len(etag)], nocache[filepath][i][1], 1)
	}
	return html
}

//保存缓存文件
func (sm *Smarty) saveCache(filepath string, html string) {
	if readcache[filepath] == false { //启用缓存&&未读取缓存则生成缓存
		//fmt.Println(nocache[filepath])
		if len(nocache[filepath]) > 0 {
			html = sm.nocacheRollback(filepath, html)
		}
		//添加删除多余空行函数
		//fmt.Println(html)
		err := ioutil.WriteFile(cachepath[filepath], []byte(html), 0666) //清空式写入
		if err != nil {
			panic("缓存写入失败" + fmt.Sprintf("%s", err))
		}
		//fmt.Println("更新了缓存" + filepath)
	}
}
