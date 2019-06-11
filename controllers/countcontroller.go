package controllers

import (
	"blogServer/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strings"
	"time"
)

/**
todo 添加接口访问计数
并通过这个controller 返回对应的计数

现在存在一个问题,有的路径是/paht/{argu}的格式,
如何处理这样的路径
*/
type siteInfo struct {
	path   string
	hasArg bool
	urlEnd int //标记在哪一个字符前的子串可以用于匹配
}

var (
	siteMap = make(map[string]*siteInfo)
)

type CountController struct {
	beego.Controller
}

//@Param tags query []int64 true "Specific tags"
//@Success 200 {object} models.QueryResult
//@router /getAccessCount [get]
func (c *CountController) AccessCount() {

}

func FilterAccessCount() {
	beego.InsertFilter("*", beego.AfterExec, func(ctx *context.Context) {
		raw, ok := accessSite("")
		if ok {
			timeStamp := time.Now().Unix()
			_ = models.AccessSite(raw, timeStamp)
		}
	})
}

//暂时没有一个合理的方式  移除查询字符串 移除path 参数
func accessSite(path string) (string, bool) {
	tmp := strings.TrimRight(path, "/")
	for k, v := range siteMap {
		if !v.hasArg {
			if tmp == v.path {
				return k, true
			}
		} else {
			chars := []rune(tmp)
			chars = chars[:v.urlEnd]
			str := string(chars)
			if str == v.path {
				return v.path, true
			}
		}
	}

	return "", false
}

func SetSite(path string, n int, pathArg bool) {
	siteMap[path] = &siteInfo{
		path:   path,
		urlEnd: n,
		hasArg: pathArg,
	}
}
