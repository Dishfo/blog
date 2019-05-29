package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"time"
)

/**
提供某一个ip访问频率的限制
*/
const (
	ipPrefix = "ip:"
	ipSet    = "ips"
)

var (
	maxAccess = 0
)

func loadConfig() {
	maxAccess = beego.AppConfig.DefaultInt("accesslimit", 60)
	maxCacheNum = beego.AppConfig.DefaultInt("maxcahcenum", 50)
	maxTopArticle = beego.AppConfig.DefaultInt("maxtoparticles", 80)
}

//UserAccessSite return true if can access this time
//return false if has over limit
func UserAccessSite(ip string) bool {
	now := time.Now()
	stamp := now.UnixNano()
	lastMin := stamp - int64(time.Minute)
	key := ipKey(ip)
	conn := client.Get()
	_, err := conn.Do("SADD", ipSet, ip)
	_, err = conn.Do("ZADD", key, stamp, stamp)
	if err != nil {
		return true
	}
	cnt, err := redis.Int(conn.Do("ZCOUNT", key, lastMin, stamp))
	if err != nil {
		return true
	}
	beego.BeeLogger.Info("%d %d", cnt, maxAccess)
	if cnt >= maxAccess {
		return false
	}
	return true
}

func ipKey(ip string) string {
	return fmt.Sprintf("%s%s", ipPrefix, ip)
}

//用于清除一分钟以前的访问记录
func clearAccessRecord() {
	now := time.Now()
	stamp := now.UnixNano()
	lastMin := stamp - int64(time.Minute)
	conn := client.Get()
	ips, err := redis.Strings(conn.Do("SMEMBERS", ipSet))
	if err != nil {
		return
	}
	for _, ip := range ips {
		key := ipKey(ip)
		_ = conn.Send("ZREMRANGEBYSCORE", key, 0, lastMin)
	}
	_ = conn.Flush()
}
